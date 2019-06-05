## filemonitor

filemonitor takes care of the storage of the files you want to monitor. It uses Git to save the files and to keep track of them. A git service, such as GitHub, could take care of notifying you when changes are detected.

Filemonitor uses the modified worker pool of [subfinder](https://github.com/subfinder/subfinder) (originally from https://github.com/stefantalpalaru/pool). The web pages/files are saved as `.txt` files, although we would love to have this changed in the near future.

Contributions are appreciated! The structure/performance has to be improved, since it's a bit unorganized/messy now. If you would like to help out, create an issue/pull request with your suggestions.

## Installation

1. Install the package.
```
go get github.com/kapytein/filemonitor
```

2. Create an empty Git repository.

2. You have to set two environment variables:

`FILEMONITOR:` Path of the Git repository you have created

`GIT_TOKEN:` filemonitor currently uses personal access tokens for authentication (might have to change this in the future). GitHub, GitLab and Bitbucket provide such tokens. It's recommended to have them as restrictive as possible.

3. The available flags:

```
FILEMONITOR v0.0.1 - Monitoring files at your wish.
Usage of ./lol:
  -beautify
        Use this if you would like to have the file 'JS beautified' when saving.
  -fetch
        Use this option if you want to start fetching the files.
  -pattern string
        This is the pattern where we will look for on the web page (in case of a dynamic link)
  -threads int
        If you choose to fetch, this is the amount of threads (concurrency) which will be used. Default is 5. (default 5)
  -url string
        This is the URL of the webpage to track (or the webpage the link is on if it is a dynamic link)
```

## How to add a URL?
If you provide a pattern, it will search for that pattern in the `src` attribute of `script` elements on the specified webpage. Other elements are not supported yet. You can create an issue if you have suggestions for different elements.

The `beautify` flag will beautify the contents upon saving as well.
`./filemonitor -url "https://google.nl" -pattern "/assets/application/main.*.js" -beautify`

## How do I specify a pattern?

Example:

`<script src="/assets/application/main.7815696ecbf1c96e6894b779456d330e.js"></script>`

Replace the dynamic value with a `*`. In above case, you would:

`-pattern "/assets/application/main.*.js"`

## How to start monitoring the web pages?

> NOTE! Do not push/change anything in your local Git repository. Let filemonitor take care of your Git repository.

You have to create cronjobs for that. Filemonitor only helps you with the storage of the files. In order to start fetching the webpages, you run the following:

`./filemonitor -fetch true -threads 20`

The tracked URLs are saved in `urls.json`, in the cwd. Once you have `cron` set up (to fetch every X hours/days), you can use the push/commit notification service of GitHub to be notified of changes. Filemonitor commits any changes found in webpages/javascript files.

#### To-Do
1. Use interfaces
2. Implement own worker pool
3. Improve naming of files
3. ??? [submit an issue! :)](https://github.com/kapytein/filemonitor/issues)
