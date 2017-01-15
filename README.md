## itunes-server

Control iTunes over a server. Uses [`osascript`](https://developer.apple.com/legacy/library/documentation/Darwin/Reference/ManPages/man1/osascript.1.html).

### Installation / setup

  - Use the following to install and run (assuming you have Go [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing)):

    ```sh
    $ go get github.com/kshvmdn/itunes-server
    $ itunes-server # Listening at localhost:8080
    ```

  - You can also build directly from source, if you'd prefer that.

    ```sh
    $ git clone https://github.com/kshvmdn/itunes-server.git
    $ cd itunes-server
    $ go build itunes-server.go
    $ ./itunes-server # Listening at localhost:8080
    ```

### Usage

  - `itunes-server` allows you to control iTunes via standard HTTP requests.
  
  - Endpoints:

    + __GET /__
      * View current status and song (if playing).
    + __GET /open__
      * Open iTunes.
    + __GET /exit__
      * Exit iTunes.
    + __GET /play__
      * Play the last-played song (iff nothing is currently playing).
    + __GET /pause__
      * Pause the currently playing song.
    + __GET /stop__
      * Stop the currently playing song.
    + __GET /next__
      * Skip to the next song.
    + __GET /prev__
      * Play the previous song.
    + __GET /mute__
      * Mute iTunes (doesn't effect system volume).
    + __GET /unmute__
      * Unmute iTunes (doesn't effect system volume).
    + __GET /shuffle__
      * Play a random song from your iTunes library.
    + __GET /tracks__
      * View list of tracks in your iTunes library.
      * Query parameters:
        * `limit` - The number of tracks to show (defaults to 100)
        * `skip` - The number of tracks to skip (defaults to 0)
    + __GET /play/track/:track_name__
      * Play track(s) with name that matches `track_name`.
    + __GET /play/artist/:artist_name__
      * Play track(s) with artist that matches `artist_name`.
    + __GET /play/album/:album_name__
      * Play track(s) with album that matches `album_name`.

  - Examples (these endpoints can be accessed through the browser as well):

    + View current status:
      
      ```sh
      curl -L localhost:8080
      {"status":"playing","current":{"title":"Gobstopper","artist":"J Dilla","album":"Donuts"}}
      ```

    + Next:

      ```sh
      curl -L localhost:8080/next
      {"status":"playing","current":{"title":"One For Ghost","artist":"J Dilla","album":"Donuts"}}
      ```

    + Pause:

      ```sh
      curl -L localhost:8080/pause
      {"status":"paused","current":{"title":"","artist":"","album":""}}
      ```

    + Play
    
      ```sh
      curl -L localhost:8080/play
      {"status":"playing","current":{"title":"One For Ghost","artist":"J Dilla","album":"Donuts"}}
      ```

  - If you're looking to expose your local server (so people can access iTunes without having to be on the same network), I suggest using [ngrok](https://ngrok.com/).

### Contribute

This project is completely open source, feel free to open an issue or create a pull request.
