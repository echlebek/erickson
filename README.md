# erickson
[![Build Status](https://api.travis-ci.org/echlebek/erickson.svg)](https://api.travis-ci.org/echlebek/erickson)

![Screenshot](/../screenshots/screenshots/screenshot_1.png?raw=true "Annotating a review")

Erickson is a simple code review app. Code reviews are created via
a command-line client, or the web frontend.

Rationale
---------
I wanted to show what could be accomplished in web-programming on the server side
with Go, with minimal dependencies and complexity. Erickson is both a functional
application, and a testbed for exploring my own ideals in software development.

Submitting a code review to erickson with git
---------------------------------------------
From within the repository you wish you review code:

    git erickson post HEAD^...

The arguments to `git erickson post` are the same arguments you'd supply to `git diff`.

Working with code reviews
-------------------------
When a review is created, the command line app returns a link to the review.
Users can annotate the review, and the owner can mark it as submitted or discarded.

Notes
-----
Erickson is pre-alpha software and may not be suitable for running in production.
Although it implements TLS, CSRF protection and secure sessions, these features
have not been vetted by a security professional. Use at your own risk.

Project goals
-------------
* Dead-simple setup that doesn't require integration with other services.
* Produces a standalone binary that doesn't require an installer or asset files.
* Simple, lightweight UI that requires minimal JS.
* Small feature-set that is robust and reliable.
* A very fast server.

Open Questions
--------------
How should erickson support notifications? E-mail? Tweets? Not sure yet.

Mercurial Support
-----------------
Dimitri Tcaciuc has written a mercurial plugin for erickson.
https://github.com/dtcaciuc/hgerickson

TODO
----
- [x] CRUD app that supports the essentials of working with code reviews.
- [x] Side-by-side diff display, rendered as HTML.
- [x] UI for annotations.
- [x] Git plugin for submitting reviews.
- [ ] Mercurial plugin for submitting reviews.
- [ ] Add support for revising reviews. (Partially done in the persistence layer)
- [ ] Remove jQuery and use plain old javascript.
- [ ] 100% test coverage.
