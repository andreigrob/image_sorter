package main

import (
	ct "context"
	"io"
	"os"
	fp "path/filepath"
	st "strings"

	dv "google.golang.org/api/drive/v3"
	ot "google.golang.org/api/option"
)

// Drive wraps the Google Drive service and destination folders.
type Drive struct {
	srv          *dv.Service
	Destinations map[string]string // map button name to folder ID
}

// NewDrive creates a Drive service using credentials from the given file.
func (d *Drive) Init(ctx ct.Context, credFile string) (e error) {
	d.srv, e = dv.NewService(ctx, ot.WithCredentialsFile(credFile))
	if e != nil {
		return
	}
	d.Destinations = make(map[string]string)
	return
}

// FindFolderID finds a folder with the given name. If multiple exist, the first is returned.
func (d *Drive) FindFolderID(name string) (_ string, e error) {
	q := "mimeType='application/vnd.google-apps.folder' and name='" + name + "'"
	res, e := d.srv.Files.List().Q(q).Fields("files(id,name)").Do()
	if e != nil {
		return
	}
	if len(res.Files) == 0 {
		return ``, os.ErrNotExist
	}
	return res.Files[0].Id, nil
}

// ListImages lists image files inside a folder.
func (d *Drive) ListImages(folderID string) (_ []*dv.File, e error) {
	q := "'" + folderID + "' in parents and (mimeType contains 'image/')"
	files := []*dv.File{}
	pageToken := ``
	var (
		res  *dv.FileList
		call *dv.FilesListCall
	)
	for {
		call = d.srv.Files.List().Q(q).Fields("nextPageToken, files(id,name)")
		if pageToken != `` {
			call.PageToken(pageToken)
		}
		res, e = call.Do()
		if e != nil {
			return
		}
		files = append(files, res.Files...)
		if res.NextPageToken == `` {
			break
		}
		pageToken = res.NextPageToken
	}
	return files, nil
}

// DownloadFile downloads a Drive file to a temporary path and returns that path.
func (d *Drive) DownloadFile(fileID, name string) (path string, e error) {
	resp, e := d.srv.Files.Get(fileID).Download()
	if e != nil {
		return
	}
	defer resp.Body.Close()

	tmpDir := os.TempDir()
	path = fp.Join(tmpDir, name)
	f, e := os.Create(path)
	if e != nil {
		return
	}
	defer f.Close()
	if _, e = io.Copy(f, resp.Body); e != nil {
		return
	}
	return
}

// MoveFile moves a Drive file to a new folder.
func (d *Drive) MoveFile(fileID, destFolderID string) (e error) {
	// Retrieve existing parents to remove
	file, e := d.srv.Files.Get(fileID).Fields("parents").Do()
	if e != nil {
		return
	}
	_, e = d.srv.Files.Update(fileID, nil).AddParents(destFolderID).RemoveParents(st.Join(file.Parents, ",")).Do()
	return
}
