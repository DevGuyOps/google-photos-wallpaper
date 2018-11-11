package main

import (
	"io/ioutil"
	"log"

	"golang.org/x/oauth2/google"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

func main() {
	b, err := ioutil.ReadFile("client_id.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, photoslibrary.PhotoslibraryReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	photoslibraryService, err := photoslibrary.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// getAlbumList(photoslibraryService)
	// getPhotosFromAlbum(photoslibraryService, "AJEwh1KefYSAZOcl9Jm7xLIwfwUyIBYCI2QVHV-vxHAK5IS0L4BoSAy4q-u8df4yW1EuPb18ciFv")
	// getPhotosFromAlbum(photoslibraryService, "AJEwh1Ks1GBG2QCp-gkji7VzHM1D1_PJlgc0g1lsNHIn0qbb6cM_Nd2vJL0ygQkoZJi5J08vkMjR")
	num := numItemsInAlbum(photoslibraryService, "AJEwh1Ks1GBG2QCp-gkji7VzHM1D1_PJlgc0g1lsNHIn0qbb6cM_Nd2vJL0ygQkoZJi5J08vkMjR")
	log.Println(num)
}

func getPhotosFromAlbum(photoslibraryService *photoslibrary.Service, albumId string) {
	searchMediaRequest := &photoslibrary.SearchMediaItemsRequest{
		AlbumId:  albumId,
		PageSize: 50,
	}

	items, err := photoslibraryService.MediaItems.Search(searchMediaRequest).Do()
	if err != nil {
		log.Println(err)
	}

	log.Println(numItemsInAlbum(photoslibraryService, albumId))

	log.Println(len(items.MediaItems))

	for _, item := range items.MediaItems {
		downloadImage("items/"+item.Id+".jpg", item.BaseUrl, item.MediaMetadata.Width, item.MediaMetadata.Height)
	}
}

func getAlbumList(photoslibraryService *photoslibrary.Service) {
	albumList, err := photoslibraryService.Albums.List().Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, album := range albumList.Albums {
		log.Println(album.Title + "-" + album.Id)
	}
}

func numItemsInAlbum(photoslibraryService *photoslibrary.Service, albumId string) int64 {
	album, err := photoslibraryService.Albums.Get(albumId).Do()
	if err != nil {
		log.Fatal(err)
	}

	return album.TotalMediaItems
}
