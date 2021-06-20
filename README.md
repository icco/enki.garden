# enki.garden

A personal internet archive.

## Background

I love the Internet Archive. I think their hoarder approach to the Internet is the only way we can preserve history for the future. We are already getting to the point where there are programs we cannot run, art we cannot see, machines we cannot fix. We need to preserve all of this so that our children's children can understand how we got to where we are today.

There is a problem though, a lot of things that we as humans acquire cannot be legally storred. Copyright law prevents us from saving and sharing our music and movies. We buy books that can only be read on a device made by the book seller, we pay someone to let us rent access to culture, and when we cannot pay anymore, our access is cut off. Also, when these businesses go out of business, we also lose access.

So until laws are fixed up, I want the ability to organize my digital hoard. I want it backed up. I want it searchable. I want it editable.

Right now I use a mixture of Github, [Kodi](http://kodi.tv/), [Sickbeard](http://sickbeard.com/), [Dropbox](http://www.dropbox.com/) and [Google Drive](https://drive.google.com) to store and organize all of my shit. I have about 20TB of Movies, TV, Random Video, Music, Photography, Text, PDFs, EPubs, Mobi docs, etc.

## Things to think about

 - Camlistore
 - https://infinit.sh/
 - Metadata storage
 - Cost
 - https://github.com/N0taN3rd/Squidwarc
 - https://github.com/internetarchive/heritrix3
 - https://archive.fo/

Camlistore has the best/closest thoughts to what I want in a system: https://camlistore.org/docs/overview

## MVP

The MVP for this app is a simple indexing and search app. It's not actually that simple though because of a few things:

 - My personal laptop
   - 116,241 dirs and 538,962 files in its home directory
   - 511,158 directories and 2,667,402 files.
   - Walking `/` with `tree -if` takes greater than ten minutes (aprox 30 min)
 - I have eight computers I want to index
   - 3 rPi
   - 2 Synology (1812+, 1813+)
     - Note: [Compile Golang for synology](http://www.faun.me/2015/02/17/cross-compiling-golang-applications-for-synology-1513.html)
   - 3 OS X boxen

So, that means worse case I probably have 20 million files on all of my computers... which is a non-trivial amount to store in a database.
