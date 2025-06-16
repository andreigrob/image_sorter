package gdrive

import (
	ls "container/list"
	ct "context"
	"fmt"
	"io"
	"log"
	"os"
	fp "path/filepath"
	st "strings"

	gc "github.com/andreigrob/image_sorter/gclient"
	tp "github.com/andreigrob/image_sorter/types"
	dv "google.golang.org/api/drive/v3"
	op "google.golang.org/api/option"
)

type DataGDriveT struct {
	gDrive *dv.Service
	Ids    tp.IdMapT
}

type GDriveT struct {
	*DataGDriveT
}

// New creates a Drive service using credentials from the given file.
func New(ctx ct.Context, credentials string) (g GDriveT, e error) {
	g.DataGDriveT = &DataGDriveT{
		Ids: make(tp.IdMapT, 20),
	}

	const tokenName = `token.json`
	const scope = dv.DriveScope

	gClient := gc.NewGClient(ctx, credentials, tokenName, scope)

	_ = gClient.CheckScope()
	httpCl := gClient.NewHttpClient()

	g.gDrive, e = dv.NewService(ctx, op.WithHTTPClient(httpCl))
	if e != nil {
		log.Printf("Unable to retrieve GDrive client: %v", e)
		return
	}

	const fields = `nextPageToken, files(id, name)`
	res, e := g.gDrive.Files.List().PageSize(100).Fields(fields).Do()
	if e != nil {
		log.Printf("Unable to retrieve GDrive files: %v", e)
		return
	}

	Len := len(res.Files)
	if Len == 0 {
		fmt.Println("No files found.")
		return
	}

	fmt.Println("Files:")
	var i int
	for ; i < Len; i++ {
		fmt.Printf("%s (%s)\n", res.Files[i].Name, res.Files[i].Id)
	}

	return
}

func (g GDriveT) getCall(query string) (_ *dv.FilesListCall) {
	return g.gDrive.Files.List().PageSize(100).Q(query)
}

// FindFolderId finds a folder with the given name. If multiple exist, the first is returned.
func (g GDriveT) FindFolderId(name tp.NameT) (_ tp.IdT, e error) {
	const query = `mimeType='application/vnd.google-apps.folder' and name='`
	list, e := g.getCall(query + string(name) + `'`).Fields(`files(id,name)`).Do()
	if e != nil {
		log.Printf("Error finding folder %s: %v", name, e)
		return
	}
	if len(list.Files) == 0 {
		log.Printf("Folder %s not found", name)
		return ``, os.ErrNotExist
	}
	log.Printf("Found %d folders: %v\n", len(list.Files), list.Files)
	return tp.IdT(list.Files[0].Id), e
}

// ListImages lists image files inside a folder.
func (g GDriveT) ListImages(folder tp.GFileT) (files []*dv.File, e error) {
	var (
		query = `'` + string(folder.Id) + `' in parents and (mimeType contains 'image/')`
		list  *dv.FileList
		call  *dv.FilesListCall
		l     = ls.New()
		i     int
	)
	for {
		call = g.getCall(query).Fields(`nextPageToken, files(id,name)`)
		if list != nil {
			call.PageToken(list.NextPageToken)
		}
		list, e = call.Do()
		if e != nil {
			log.Printf("Error listing images in %v: %v", folder, e)
			return
		}

		for ; i < len(list.Files); i++ {
			l.PushBack(list.Files[i])
		}
		if list.NextPageToken == `` {
			break
		}
	}

	files = make([]*dv.File, l.Len())
	i = 0
	var v *ls.Element = l.Front()
	for ; v != nil; v = v.Next() {
		files[i] = v.Value.(*dv.File)
		i++
	}

	return
}

// Download downloads a Drive file to a temporary path and returns that path.
func (g GDriveT) Download(file tp.GFileT) (path string, e error) {
	res, e := g.gDrive.Files.Get(string(file.Id)).Download()
	if e != nil {
		log.Printf("Error downloading %v: %v", file, e)
		return
	}
	defer res.Body.Close()

	path = fp.Join(os.TempDir(), string(file.Name))
	f, e := os.Create(path)
	if e != nil {
		log.Printf("Error creating %s: %v", path, e)
		return
	}
	defer f.Close()
	if _, e = io.Copy(f, res.Body); e != nil {
		log.Printf("Error downloading %s to %s: %v", file.Name, path, e)
		return
	}

	log.Printf("Downloaded %s to %s", file.Name, path)
	return
}

// Move moves a Drive file to a new folder.
func (g GDriveT) Move(file, folder tp.GFileT) (e error) {
	// Retrieve existing parents to remove
	fileData, e := g.gDrive.Files.Get(string(file.Id)).Fields(`parents`).Do()
	if e != nil {
		log.Printf("Error retrieving parent folders of %v: %v", file, e)
		return
	}

	log.Printf("Parent folders of %v: %v", file, fileData.Parents)
	_, e = g.gDrive.Files.Update(string(file.Id), nil).AddParents(string(folder.Id)).RemoveParents(st.Join(fileData.Parents, `,`)).Do()
	if e != nil {
		log.Printf("Error moving %v to %v: %v", file, folder, e)
		return
	}

	log.Printf("Moved %v to %v", file, folder)
	return
}
