package main

import (
    "fmt"
    "github.com/andygrunwald/go-jira"
    "github.com/dustin/go-humanize"
    "github.com/mitchellh/go-homedir"
    "gopkg.in/cheggaaa/pb.v1"
    "flag"
    "os"
    "io"
    "io/ioutil"
    "encoding/json"
    "bufio"
)

type JiraInstance struct {
    Name string `json:"name"`
    Url string `json:"url"`
}

type Config struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Sites []JiraInstance `json:"sites"`
}

func download(c *jira.Client, a *jira.Attachment) error {

    if _, err := os.Stat(a.Filename); err == nil {
        fmt.Printf("[SKIP] %s\n", a.Filename)
    } else {

        output, err := os.Create(a.Filename)

        if err != nil {
            return err
        }

        defer output.Close()

        human_size := humanize.Bytes(uint64(a.Size))

        fmt.Printf("[DOWNLOAD] %s (%s)\n", a.Filename, human_size)

        bar := pb.New(a.Size)
        bar.SetUnits(pb.U_BYTES)
        bar.SetMaxWidth(80)
        bar.Start()

        resp, err := c.Issue.DownloadAttachment(a.ID)

        if err != nil {
            return err
        }

        reader := bar.NewProxyReader(resp.Response.Body)

        io.Copy(output, reader)

        bar.Finish()

    }

    return nil
}


func download_all(c *jira.Client, i *jira.Issue) error {

   for _, attachment := range i.Fields.Attachments {
       err := download(c, attachment)

       if err != nil {
           return err
       }
   }

   return nil
}

func upload(c *jira.Client, i *jira.Issue, filename string) error {
    fmt.Printf("I should upload %s\n", filename)

    stat, err := os.Stat(filename)

    if err != nil {
        return err
    }

    human_size := humanize.Bytes(uint64(stat.Size()))

    fmt.Printf("[UPLOAD] %s (%s)\n", filename, human_size)

    fh, _ := os.Open(filename)

    bar := pb.New(int(stat.Size()))
    bar.SetUnits(pb.U_BYTES)
    bar.SetMaxWidth(80)
    bar.Start()

    reader := bar.NewProxyReader(bufio.NewReader(fh))

    _, _, err =c.Issue.PostAttachment(i.ID, reader, filename)

    if err != nil {
        panic(err)
    }

    bar.Finish()
    return nil
}

func read_config(path string) (c Config, err error){
    data, err := ioutil.ReadFile(path)

    if err != nil {
        return c, err
    }

    err = json.Unmarshal(data, &c)

    if err != nil {
        return c, err
    }

    return c, nil
}

func main() {

    var issue_ref = flag.String("issue", "", "Issue references")
    var loc = flag.String("loc", "", "Site to use as defined in config")
    var u_file = flag.String("upload", "", "File to upload")

    flag.Parse()

    if *issue_ref == "" {
        fmt.Printf("Please specify a -issue to sync.\n")
        os.Exit(2)
    }

    if *loc == "" {
        fmt.Printf("Please specify a -loc to use.\n")
        os.Exit(2)
    }

    // Read config
    config_path, err := homedir.Expand("~/.jira-file-sync.json")


    if err != nil {
        panic(err)
    }

    c, err := read_config(config_path)

    if err != nil {
        panic(err)
    }


    // Determine URL to use
    var url = ""

    for _, site := range c.Sites {

        if *loc == site.Name {
            url = site.Url
        }
    }

    if url == "" {
        panic("Did not set URL")
    }

    // Initiate Client
    client, err := jira.NewClient(nil, url)
    if err != nil {
        panic(err)
    }

    res, err := client.Authentication.AcquireSessionCookie(c.Username, c.Password)
    if err != nil || res == false {
        fmt.Printf("Result: %v\n", res)
        panic(err)
    }


    issue, _, err := client.Issue.Get(*issue_ref)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)


    if *u_file != "" {
        err = upload(client, issue, *u_file)
    } else {
        err = download_all(client, issue)
    }

    if err != nil {
        panic(err)
    }

}
