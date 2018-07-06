# Core
This is the core package that should reside on the machine of the person using Hackerlog. It will
communicate to the Hackerlog API.

## Command line args
```
usage: hackerlog [-h|--help] -u|--api-url "<value>" -t|--editor-token "<value>"
                 -e|--editor-type "<value>" -p|--project-name "<value>"
                 -f|--file-name "<value>" [-w|--loc-written <integer>]
                 [-d|--loc-deleted <integer>] -s|--started-at "<value>"
                 -x|--stopped-at "<value>"

                 Collects coding stats and submits them to the API.

Arguments:

  -h  --help          Print help information
  -u  --api-url       The URL of the API to send the request
  -t  --editor-token  The editor token associated with a user
  -e  --editor-type   The editor that is being used
  -p  --project-name  The name of the project associated with the unit of work.
  -f  --file-name     The file name that was edited
  -w  --loc-written   The amount of lines of code that has been written.
  -d  --loc-deleted   The amount of lines of code that has been deleted.
  -s  --started-at    When did the file start being edited
  -x  --stopped-at    When did the file stop being edited
```

```
./core -u "http://localhost:8080/v1/units" -t "8ccb0396-1f14-49f2-b050-fd9208b8f6b2" -e "vscode" -p "core" -f "server.go" -w 224 -d 12 -s "2018-07-06T14:16:43+00:00" -x "2018-07-06T16:16:43+00:00"
```
