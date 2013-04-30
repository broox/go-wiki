# Gowiki

Playing around with Go by doing a simple wiki tutorial and extending upon it a bit.

## Setup

1. Setup Go Project directories

        /your-go-path/bin
        /your-go-path/pkg
        /your-go-path/src

2. Define your path variables in `~/.profile`

        export GOPATH=/your-go-path
        export PATH=$PATH:$HOME/your-go-path/bin

3. Restart bash or run `source ~.profile`

4. Clone gowiki src into `/your-go-path/src` directory

        git clone git@github.com:broox/go-wiki.git

5. Create MySQL DB called `gowiki` with the following schema.

        create table `pages` (
            `id` int(11) NOT NULL AUTO_INCREMENT,
            `title` varchar(100),
            `body` text,
            `created_at` datetime,
            `updated_at` datetime,
            primary key (`id`)
        );

6. Build the project via `go build gowiki`

7. Run the application via `gowiki`
