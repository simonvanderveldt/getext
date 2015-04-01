package main

// From https://www.socketloop.com/tutorials/golang-untar-or-extract-tar-ball-archive-example
// and http://www.qetee.com/development/golang-tar.html

import (
    "archive/tar"
    "compress/gzip"
    "flag"
    "fmt"
    "io"
    "os"
    "strings"
    "path"
    "net/http"
)

func main() {
    // Get the command-line arguments
    flag.Parse() 

    // Exit if not exactly 1 argument is given
    if len(flag.Args()) == 0 || len(flag.Args()) > 1 { 
        Usage()
        os.Exit(1)
    }

    destDirPath := "."

    // Get the source URL from the first and only argument
    sourceurl := flag.Arg(0)

    // Check if the sourceurl starts with http(s)://
    if strings.HasPrefix(sourceurl, "http://") != true && strings.HasPrefix(sourceurl, "https://") != true {
        Usage()
        os.Exit(1)
    }
    // Check if the sourceurl ends with .tar.gz
    if strings.HasSuffix(sourceurl, ".tar.gz") != true {
        Usage()
        os.Exit(1)
    }

    // Open the url
    response, err := http.Get(sourceurl)
    if err != nil {
        fmt.Println("Error while downloading", sourceurl, "-", err)
        return
    }
    defer response.Body.Close()

    // Read the file as gzipped files
    gzipReader, err := gzip.NewReader(response.Body)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer gzipReader.Close()

    // Now read the decompressed file as a tarball
    tarBallReader := tar.NewReader(gzipReader)
    
    // Extracting tarred files
    for {
        // Advance to the next file in the tarball
        header, err := tarBallReader.Next()
        if err != nil {
            if err == io.EOF {
                // End of tarball, exit
                break
            }
            fmt.Println(err)
            os.Exit(1)
        }

        // Skip pax_global_header with the commit ID this archive was created from
        if header.Name == "pax_global_header" {
            continue
        }

        // Now extract the current file
        fmt.Println("Extracting file " + header.Name)
        // Check if the current file is not a directory (means we won't extract empty directory's)
        if header.Typeflag != tar.TypeDir {
            
            // Create the directory before we create the file
            os.MkdirAll(destDirPath + "/" + path.Dir(header.Name), os.ModePerm)
            
            // Create the file to extract to
            fileWriter, err := os.Create(destDirPath + "/" + header.Name)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            
            // Copy the contents of the current file from the tarball to the target
            _, err = io.Copy(fileWriter, tarBallReader)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            err = os.Chmod(header.Name, os.FileMode(header.Mode))
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
    }
}

var Usage = func() {
        fmt.Println("Usage: getext sourceurl")
        fmt.Println("sourceurl should be of the format http://*.tar.gz")
}