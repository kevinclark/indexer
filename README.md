This is a very basic search-engine engine that I'm writing to help me
learn more about go.

If you'd like to play with it, `index.go` at top level is an example
entry point that will index a maildir path and is able to load the index
back in for querying. Currently it doesn't support snippets and the
inverted index doesn't have term locations - just existence in a file.

This was developed against the [enron corpus](https://www.cs.cmu.edu/~enron/),
but since it's using go's built in mail parser, I expect it should work
for maildir's in general.

It's possible you shouldn't copy idioms in this repo since I'm just learning
them myself. This might turn into something real, but isn't currently
intended for anything in production.

So take everything here with a grain of salt.
