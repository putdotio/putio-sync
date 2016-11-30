// Package putio is the Put.io API v2 client for Go.
//
// The go-putio package does not directly handle authentication. Instead, when
// creating a new client, pass an http.Client that can handle authentication for
// you. The easiest and recommended way to do this is using the golang.org/x/oauth2
// library, but you can always use any other library that provides an http.Client.
// If you have an OAuth2 access token (for example, a personal API token), you can
// use it with the oauth2 library using:
// 	import "golang.org/x/oauth2"
//
// 	func main() {
// 		tokenSource := oauth2.StaticTokenSource(
// 			&oauth2.Token{AccessToken: "<YOUR-TOKEN-HERE>"},
// 		)
// 		oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
// 		client := putio.NewClient(oauthClient)
// 		// get root directory
// 		root, err := client.Files.Get(0)
// 	}
// Note that when using an authenticated Client, all calls made by the client will
// include the specified OAuth token. Therefore, authenticated clients should
// almost never be shared between different users.
package putio
