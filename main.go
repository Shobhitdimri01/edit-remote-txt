package EditTxt

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"gopkg.in/mcuadros/go-defaults.v1"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type NewTxt struct{
        S1apIP string	`default:"10.55.55.1"`
        PLMNID string	`default:"20893"`
        MCC    string	`default:"208"`
}
type Remoteserverconfig struct{
        Username        string  `default:"username"`
        Password        string  `default:"user_password"`
        Serverip        string  `default:"remoteip:22"`
		Protocol		string	`default:"tcp"`
}

//*****----setting default structs----******
func SetRemoteServer() *Remoteserverconfig {
    example := new(Remoteserverconfig)
    defaults.SetDefaults(example) //<-- This set the defaults values

    return example
}
func SetMyTxt() *NewTxt{
	SetNewTxt := new(NewTxt)
	defaults.SetDefaults(SetNewTxt)
	return SetNewTxt
}
var localpath ="/home/localpath/test.txt"
var dstpath = "/home/username/testtxt/test.txt"
//****----- file editing func ------****
func WriteFile(){
	remote := SetRemoteServer()
        sshconfig := &ssh.ClientConfig{
                User: remote.Username,
                Auth: []ssh.AuthMethod{
                        ssh.Password(remote.Password),
                },
                HostKeyCallback: ssh.InsecureIgnoreHostKey(),// optional
        }

        client,err := ssh.Dial(remote.Protocol,remote.Serverip,sshconfig)
        if err != nil {
                panic("failed to dail: "+err.Error())
        }
        fmt.Println("Succesfully connected to "+remote.Serverip+" server")
	sftp,err :=sftp.NewClient(client)
        if err != nil {
                log.Fatal(err)
        }
        defer sftp.Close()
	file,err := sftp.Open(dstpath)
        if err != nil {
                fmt.Println(err)
        }
        input, err := io.ReadAll(file)
        if err != nil {
                        fmt.Println(err)
        }
	
        lines := strings.Split(string(input), "\n")

		newtxt :=SetMyTxt()
		//******----setting text to be edit-----******
        Edit(lines,"s1ap",newtxt.S1apIP)
        Edit(lines,"plmn_id",newtxt.PLMNID)
        Edit(lines,"mcc",newtxt.MCC)
	  output := strings.Join(lines, "\n")
                
                ip:=strings.Split(string(remote.Serverip), ":")
                remotepath :=remote.Username+"@"+ip[0]+":"+dstpath
                err = ioutil.WriteFile(localpath, []byte(output), 0644)
        if err != nil {
               fmt.Println(err)
        }

		//****---copying file from local to remote server----****
                copyfile := "sshpass -p "+remote.Password+" scp -r "+localpath+" "+remotepath
                cmd := exec.Command("bash","-c",copyfile)
                fmt.Println("Remote file edited")
                stdout,err := cmd.Output()
                if err != nil {
                        fmt.Println(err)
                }
                //printing output
                fmt.Println(string(stdout))
}

//****---main editing logic---****
func Edit(lines []string,txt string,Key string){
        for i, line := range lines {
                if strings.Contains(line, txt) {
                        lines[i] = txt+": "+Key
                }
        }
}