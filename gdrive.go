package main

import (
    "context"
    "io"
    "os"
    "path/filepath"
    "strings"

    "google.golang.org/api/drive/v3"
    "google.golang.org/api/option"
)

// Drive wraps the Google Drive service and destination folders.
type Drive struct {
    srv          *drive.Service
    Destinations map[string]string // map button name to folder ID
}

// NewDrive creates a Drive service using credentials from the given file.
func NewDrive(ctx context.Context, credFile string) (*Drive, error) {
    srv, err := drive.NewService(ctx, option.WithCredentialsFile(credFile))
    if err != nil {
        return nil, err
    }
    return &Drive{srv: srv, Destinations: make(map[string]string)}, nil
}

// FindFolderID finds a folder with the given name. If multiple exist, the first is returned.
func (d *Drive) FindFolderID(name string) (string, error) {
    q := "mimeType='application/vnd.google-apps.folder' and name='" + name + "'"
    res, err := d.srv.Files.List().Q(q).Fields("files(id,name)").Do()
    if err != nil {
        return "", err
    }
    if len(res.Files) == 0 {
        return "", os.ErrNotExist
    }
    return res.Files[0].Id, nil
}

// ListImages lists image files inside a folder.
func (d *Drive) ListImages(folderID string) ([]*drive.File, error) {
    q := "'" + folderID + "' in parents and (mimeType contains 'image/')"
    files := []*drive.File{}
    pageToken := ""
    for {
        call := d.srv.Files.List().Q(q).Fields("nextPageToken, files(id,name)")
        if pageToken != "" {
            call.PageToken(pageToken)
        }
        res, err := call.Do()
        if err != nil {
            return nil, err
        }
        files = append(files, res.Files...)
        if res.NextPageToken == "" {
            break
        }
        pageToken = res.NextPageToken
    }
    return files, nil
}

// DownloadFile downloads a Drive file to a temporary path and returns that path.
func (d *Drive) DownloadFile(fileID, name string) (string, error) {
    resp, err := d.srv.Files.Get(fileID).Download()
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    tmpDir := os.TempDir()
    path := filepath.Join(tmpDir, name)
    f, err := os.Create(path)
    if err != nil {
        return "", err
    }
    defer f.Close()
    if _, err := io.Copy(f, resp.Body); err != nil {
        return "", err
    }
    return path, nil
}

// MoveFile moves a Drive file to a new folder.
func (d *Drive) MoveFile(fileID, destFolderID string) error {
    // Retrieve existing parents to remove
    file, err := d.srv.Files.Get(fileID).Fields("parents").Do()
    if err != nil {
        return err
    }
    _, err = d.srv.Files.Update(fileID, nil).
        AddParents(destFolderID).
        RemoveParents(strings.Join(file.Parents, ",")).
        Do()
    return err
}
