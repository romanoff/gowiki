Gowiki
======

Gowiki is a project that should make creation of websites with documentation easier. All documentation pages are created in markdown format. No configuration required. There is a concept of project, section and article. Those are based on files/folders structure. Here is simple example:

```
gowiki/
├── 1_getting_started
└── 2_usage
```

In this example we have one folder and 2 articles. Project name will be `Gowiki` and article names will be `Getting started` and `Usage`. Number prefixes in article file names are used for order (they are optional). No sections are present in this example. Folders inside a project would stand for sections. Example:

```
interesting_projects/
├── datakeeper
│   └── getting_started
└── gowiki
    ├── 1_getting_started
    └── 2_usage
```

Here we have `Interesting projects` project and two sections under it - `Datakeeper` and `Gowiki`. Each project has articles.

Sections are represented with icon and text. Material design icons are included into gowiki binary. To specify what icon do you want to use in section, you can add `__<icon_name>` to folder name. Icons can be found [here](https://www.google.com/design/icons/) (just replace spaces in names with underscores). 

You can create as many subsections as you want. Also by using `gowiki`, you get full text search powered by [bleve](https://github.com/blevesearch/bleve).

Binary releases can be downloaded [here](https://github.com/romanoff/gowiki/releases) or compiled from source.

Starting gowiki
===============

To start web server serving you wiki content, just go to the root of your project `gowiki` project and execute `gowiki` command in the terminal. After this your wiki should be served from `http://127.0.0.1:8080/`.
