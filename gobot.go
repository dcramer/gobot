package main

import (
    "bytes"
    "fmt"
    "github.com/lye/cleanirc"
    "log"
    "os"
    "os/exec"
    "regexp"
    "strings"
)

const IRC_NICK = "gobot"
const IRC_SERVER = "irc.freenode.net"

func main() {
    // TODO: we should recompile this anytime the bot's nickname changes
    msgparser, error := regexp.Compile(fmt.Sprintf("^%s[:,] (?P<message>.+)$", IRC_NICK))
    if error != nil {
        log.Fatal(error)
    }

    irccon := irc.IRC(IRC_NICK, IRC_NICK)
    irccon.VerboseCallbackHandler = false
    
    error = irccon.Connect(IRC_SERVER)
    if error != nil {
        fmt.Printf("%s\n", error)
        fmt.Printf("%#v\n", irccon)
        os.Exit(1)
    }
    irccon.AddCallback("001", func(e *irc.IRCEvent) {
        // irccon.Join("#disqus")
        // irccon.Join("#ops")
    })
    irccon.AddCallback("PRIVMSG", func(e *irc.IRCEvent) { 
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

func handleMessage(message string, event *irc.IRCEvent, respond func(output string)) {
    switch message {
        // case "what's deployed?":
        //     latestRevision := getWhatsDeployed()
        //     if latestRevision != "" {
        //         respond(fmt.Sprintf("https://github.com/dcramer/gobot/commits/%s", latestRevision))
        //     }
    }
}