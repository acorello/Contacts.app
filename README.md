# A REST-ful DDD-organized Web app to learn Go, its stdlib, and HTMX

The application replicates exactly the one used in the book [Hypermedia Systems](https://hypermedia.systems) (a tutorial about HTMX, HyperView and REST-fulness) but it's designed according to some self-imposed constraints and design principles.

## Project Goals

- review REST-ful (as Fielding's dissertation) architecture
- learn fundamentals of Go Web development
- learn fundamentals of HTMX
- learn Go standard library
- sharpen my Go programming skills
- draft the way I think a web-app and a project should be organized according to my take on DDD

  The tutorial presents a single entity: the Contact. But I'm writing the project as if more will come.

## Project Non-Goals

- be a production-realistic example
- authentication and security features
- data-persistance
- DRY-out all the things
- …etc

## Constraints

1. use only Go stdlib and don't spend time re-inventing a fancy web framework

   1. therefore, I haven't used paths pattern-matching (eg. `/myentity/:id/property`). Instead I pass all the dynamic values in the URL query
   1. no reflection or generic logic to associate http-handlers with http methods. Instead I use a `switch-case` on the http-method value.

1. no hardcoded application URLs in the HTML templates…

   all URL's are passed (pre-computed if dynamic) as template parameters. Also see the [Design Principles](#design-principles)

1. only implement the functionality presented in the tutorial or less (but I may implement it differently, eg. by using a different HTTP method)

## Shortcuts

1. only an in-memory db, but design program against an interface (not a concrete implementation)
1. the DB interface stands for a component performing I/O and so should accept a `context` and return an error in all methods; I haven't bothered because of the [Project Non-Goals](#project-non-goals)
1. no tests (unless for exploratory reasons)

   the goal here is to learn not to deliver a production system; this project is a canvas where I'm painting ideas using code

1. just-enough CSS

   Aesthetic is not a goal here but we also don't want our eyes to bleed; so I just styled it with [PicoCSS](https://picocss.com) with default settings… and semantic HTML is all I need write, sweet.

1. tidy-up templates setup later

   I'm not happy with how I'm parsing and loading templates but I haven't found way that seems idiomatic, logical, and is optimal (parse each template only once and then compose the parsed-trees)

## Design Principles

I'm not implementing a full-featured app, I'm just implementing what the tutorial presented my way, with an eye to the future, so the design principles listed here are just the ones I consciously decided and implemented up to the point I found the tutorial relevant for my goals.

- application code is grouped by entity

  all code that represents that entity lives within that entity folder (incl. http handlers, html templates, etc.)

- make as much of the code URL agnostic as possible

  a REST-ful principle is that clients should be agnostic of the URL structure. Strictly speaking, the http-handler also doesn't need to know the URL path and I decided to try pushing this further and see what code end up writing and what architectural properties I would get if I made as much of the components of the app URL agnostic.

  I'm not sure if I will get valuable properties from this principles but I took it as an exercise in single-responsibility principles, contract design between the http handler and its environment; perhaps implementing configurable or dynamic paths could be useful although http-headers and query parameters are probably enough for any use case. Just consider it a 'calisthenics' exercise, if you don't see the point.

  One result you will notice is that the handler declares the kinds of URLs it supports for its entity (one for viewing, one for listing, one for creating a new instance, etc.) and it's allowed to add parameters (eg. the resource ID), the HTTP mux decides the actual URL root-path to assign to an HTTP handler.
