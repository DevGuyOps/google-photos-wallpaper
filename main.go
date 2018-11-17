package main

import (
	"io/ioutil"
	"log"
	"strconv"

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
	// num := numItemsInAlbum(photoslibraryService, "AJEwh1Ks1GBG2QCp-gkji7VzHM1D1_PJlgc0g1lsNHIn0qbb6cM_Nd2vJL0ygQkoZJi5J08vkMjR")
	// log.Println(num)

	getRandomPhotoFromAlbum(photoslibraryService, "AJEwh1Ks1GBG2QCp-gkji7VzHM1D1_PJlgc0g1lsNHIn0qbb6cM_Nd2vJL0ygQkoZJi5J08vkMjR")
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

	numInAlbum := numItemsInAlbum(photoslibraryService, albumId)

	log.Println("here: " + string(random(0, int(numInAlbum))))

	for _, item := range items.MediaItems {
		downloadImage("items/"+item.Id+".jpg", item.BaseUrl, item.MediaMetadata.Width, item.MediaMetadata.Height)
	}
}

func getRandomPhotoFromAlbum(photoslibraryService *photoslibrary.Service, albumId string) {
	numInAlbum := numItemsInAlbum(photoslibraryService, albumId)
	photoNumInAlbum := 80 // random(0, numInAlbum)
	pageSize := 50
	pageCount := 0
	nextPageToken := ""

	log.Println("Rand Num: " + strconv.Itoa(photoNumInAlbum))

	for pageCount*pageSize < numInAlbum {
		var searchMediaRequest photoslibrary.SearchMediaItemsRequest

		// Get next page details
		if nextPageToken != "" {
			searchMediaRequest = photoslibrary.SearchMediaItemsRequest{
				AlbumId:   albumId,
				PageSize:  int64(pageSize),
				PageToken: nextPageToken,
			}
		} else {
			searchMediaRequest = photoslibrary.SearchMediaItemsRequest{
				AlbumId:  albumId,
				PageSize: int64(pageSize),
			}
		}

		// Perform search
		items, err := photoslibraryService.MediaItems.Search(&searchMediaRequest).Do()
		if err != nil {
			log.Println(err)
		}

		nextPageToken = items.NextPageToken

		// Download pic and do many things with it
		if pageCount*pageSize < photoNumInAlbum && pageCount*pageSize+pageSize > photoNumInAlbum {
			// TODO: Download and stuff
			log.Println("PIC FOUND")
			break
		}

		log.Println("LOOP")
		log.Println(pageCount * pageSize)

		pageCount++
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

func numItemsInAlbum(photoslibraryService *photoslibrary.Service, albumId string) int {
	album, err := photoslibraryService.Albums.Get(albumId).Do()
	if err != nil {
		log.Fatal(err)
	}

	return int(album.MediaItemsCount)
}
