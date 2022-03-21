# Help stop Kremlin propaganda

Russia has attacked a democratic European country.

### Introduction

> This tool attempts to automate the search process of propaganda news.

We need tools to help fight propaganda and spread the word to the Russian population about the attack that's
been going on since the 24th of February 2022.

This repository holds a micro-service to scan the internet for latest news articles containing keywords and phrases
which are likely to be used to hide the fact that Russia is the aggressor and there is no peace-making mission happening
in Ukraine. It is an **occupation** of Ukrainian lands and people.

**Europe as well as the rest of the world are with Ukraine.**

---

This is a work in progress, things will break.

### Roadmap

Need help from someone who understands the linguistics of the Russian language.

- [ ] Support multiple search engines
    - [x] Google Search
    - [ ] Google Search via Proxies
    - [ ] [oxylabs.io](https://oxylabs.io/), [brightdata.com](https://brightdata.com), [smartproxy.com](https://smartproxy.com/) or similar.
    - [ ] Yandex
    - [ ] Hardcoded list of results
- [ ] Export results
    - [x] CSV
    - [x] JSON
    - [ ] REST API
    - [ ] Websockets
- [ ] CI
- [ ] docker-compose.yml
- [ ] Elasticsearch backend for article, keyword and metadata storage and analysis
- [ ] Build keyword lists
- [ ] Interfaces for result scoring and reporting
- [ ] API server and client with examples
- [ ] CD
- [ ] Reports
    - [ ] List of domains ordered by score
    
### Getting started

Code is written in Go and for most cases should be covered by unit tests.

The foundation is built on top of [Cobra](https://github.com/spf13/cobra). Project uses go modules.

### Running

Search and output results to file or stdout:
```shell
                                  [json|csv]      
go run . search keyword1 keyword2 -o -w 
                                     [stdout|file]
```

### Help

```shell
Search multiple search engines

Usage:
  drone search [flags]

Flags:
  -a, --append             append results to an existing file
  -d, --directory string   directory path for search results (default "./output")
  -o, --format string      output format (csv, json) (default "json")
  -h, --help               help for search
  -i, --import string      newline separated keyword file
  -f, --overwrite          overwrite an existing file
  -p, --pages int          specify how many pages to scroll through (default 1)
  -q, --quiet              disable logging
  -s, --suffix string      output file suffix {query}_{engine}_{suffix}.{ext}
  -w, --writer string      output writer (stdout, file) (default "stdout")

Global Flags:
  -l, --lang string        language for search engine configuration (default "ru")
  -v, --verbosity string   log level (trace, debug, info, warn, error, fatal, panic) (default "trace")
```

### Contribution

To say pseudo-anonymous I have created a separate Github account, you can do the same.

### Generate a private key

Store keys in subdirectory:
```shell
mkdir -p ~/.ssh/globalshield
```

Generate a private and public key pair:
```shell
ssh-keygen -t ed25519 -C "someemail@mailservice.com"
```

Type in full directory path and keyfile name:
```shell
/home/<user>/.ssh/globalshield/team-key
or on Mac
/Users/<user>/.ssh/globalshield/team-key
```

Symlink your keys to project dir:
```shell
ln -s ~/.ssh/globalshield/team-key `pwd`/gs.key
ln -s ~/.ssh/globalshield/team-key.pub `pwd`/gs.key.pub
```

| gs.key and gs.key.pub files are in .gitignore.

### Build git container

```shell
docker build -t drone-git -f Dockerfile-git .
```

Run Git command:
```shell
docker run --rm -it -v $(pwd)/gs.key:/root/.ssh/ed25519 -v $(pwd):/app drone-git <command>
```

To commit:
```shell
docker run --rm -it -v $(pwd)/gs.key:/root/.ssh/ed25519 -v $(pwd):/app drone-git commit -m 'message'
```

To push:
```shell
docker run --rm -it -v $(pwd)/gs.key:/root/.ssh/ed25519 -v $(pwd):/app drone-git push origin main
```

[comment]: <> (remote add origin git@github.com:globalshield/drone.git)