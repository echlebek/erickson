# erickson
[![Build Status](https://api.travis-ci.org/echlebek/erickson.svg)](https://api.travis-ci.org/echlebek/erickson)

![Screenshot](/../screenshots/screenshots/screenshot_1.png?raw=true "Annotating a review")

Erickson is a simple code review app. Code reviews are created via
the API methods, or the web frontend.

Once a review is created, annotations can be made to the diff, and
successive versions of the diff can be appended to the original review.

Erickson is a work in progress and many features are incomplete, missing
or broken.

Project goals:
* Dead-simple setup that doesn't require integration with other services.
* A standalone binary that doesn't require an installer or asset files.
* A simple, lightweight UI that requires minimal JS.
* A small feature-set that is robust and reliable.

TODO:
* Write a command-line tool for submitting and updating reviews
* Add more support for revisions
* Consider removing jQuery
* Improve test coverage
* SCM-specific tools
