package main

import (
    "crypto/tls"
    "fmt"
    "github.com/thoj/go-ircevent"
    "github.com/stathat/jconfig"
    "log"
    "os"
    "regexp"
)


func main() {
    config := jconfig.LoadConfig("gobot.json")

    ircserver := config.GetString("server")
    ircnick := config.GetString("nick")

    // TODO: we should recompile this anytime the bot's nickname changes
    msgparser, error := regexp.Compile(fmt.Sprintf("^%s[:,] (?P<message>.+)$", ircnick))
    if error != nil {
        log.Fatal(error)
    }

    irccon := irc.IRC(ircnick, ircnick)
    ircpass := config.GetString("password")
    if ircpass != "" {
        irccon.Password = ircpass
    }

    irccon.VerboseCallbackHandler = false
    irccon.UseTLS = config.GetBool("ssl")
    irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    error = irccon.Connect(ircserver)
    if error != nil {
        fmt.Printf("%s\n", error)
        fmt.Printf("%#v\n", irccon)
        os.Exit(1)
    }
    irccon.AddCallback("001", func(e *irc.Event) {
        for _, channel := range config.GetArray("channels") {
            irccon.Join(fmt.Sprintf("%s", channel))
        }
    })
    irccon.AddCallback("PRIVMSG", func(e *irc.Event) { 
        log.Println(e.Message)

        match := msgparser.FindSubmatch([]byte(e.Message))
        if match == nil {
            return
        }

        channel := e.Arguments[0]
        message := string(match[1])

        respond := func (output string) {
            irccon.Privmsg(channel, fmt.Sprintf("%s: %s", e.Nick, output))
        }

        handleMessage(message, e, respond)
    })
    irccon.Loop()
}

// func getWhatsDeployed() string {
//     var stdout bytes.Buffer

//     cmd := exec.Command("git", "ls-remote", "--tags", "git@github.com:dcramer/gobot.git", "deploy-production-*")
//     log.Print(cmd.Args)
//     cmd.Stdout = &stdout
//     error := cmd.Run()
//     if error != nil {
//         log.Print(error)
//     }

//     allLines := strings.Split(stdout.String(), "\n")
//     line := strings.TrimSpace(allLines[len(allLines) - 2])
//     latestRevision := strings.Split(line, "\t")[0]

//     return latestRevision
// }

func handleMessage(message string, event *irc.Event, respond func(output string)) {
    switch message {
        // case "what's deployed?":
        //     latestRevision := getWhatsDeployed()
        //     if latestRevision != "" {
        //         respond(fmt.Sprintf("https://github.com/dcramer/gobot/commits/%s", latestRevision))
        //     }
    }
}