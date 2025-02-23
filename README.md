# A REST-ful DDD-organized Web app to learn Go, its stdlib, and HTMX

## Live [Demo](https://contacts-app.fly.dev) deployed on [fly.io](https://fly.io)

The application replicates exactly the one used in the book [Hypermedia Systems](https://hypermedia.systems) (a tutorial about HTMX, HyperView and REST-fulness) but it's designed according to some self-imposed constraints and design principles.

The tutorial presents a single entity: the Contact; but I'm writing the project as if more will come.

## HTMX: an example

```html
<input
  type="email"
  name="Email"
  placeholder="Email"
  value="{{ .Email }}"
  hx-patch="/contact/email"
  hx-target="next .error" />
<span class="error"></span>
```

The `input[type=email]` becomes an HTML control that will fire a `PATCH /contact/email` on the default event `changed` (i.e. focus moved to the nest element and content changed).

The `hx-target="next .error"` instructs the browser to replace the following `span` element with the response; the `next .error` is [HyperScript](https://hyperscript.org) code.

More examples and excellent docs about HTMX can be found on [the HTMX docs homepage](https://htmx.org/docs/).

## Application and HTMX features implemented

The tutorial uses the requirements of an address-book CRUD application to demonstrate how to:

- remove the jarring "full-page refresh" for a smooth UI\UX (no flash of unstyled content) without changing the server-side implemntation, by just adding the `hx-boost` to the HTML body
- implement a simple search functionality
- implement pagination with on-demand or continuous scrolling
- send request with HTTP methods not natively supported by HTML (eg. DELETE)
- implement server-side valiadation of an individual fields and display of validation response

## My Goals

- see how HTMX facilitates a REST-ful (as Fielding's dissertation) architecture
- practice fundamentals of web development in Go and its standard library
- learn fundamentals of HTMX
- draft a DDD-inspired project structure

## My Non-Goals

- be a production-realistic example (eg. don't worry about authentication, security, observability, testing, etc.)
- precise validation and error reporting
- data-persistance
- …etc

## Constraints

1. use only Go stdlib

   1. therefore, I haven't used paths pattern-matching (eg. `/myentity/:id/property`). Instead I pass all the dynamic values in the URL query
   1. no reflection or generic logic to associate http-handlers with http methods. Instead I use a `switch-case` on the http-method value.

1. no hardcoded application URLs in the HTML templates…

   all URL's are passed (pre-computed if dynamic) as template parameters. Also see the [Design Principles](#design-principles)

1. only implement the functionality presented in the tutorial or less (but I may implement it differently, eg. by using a different HTTP method)

## Shortcuts

1. only an in-memory db, but design program against an interface (not a concrete implementation)
1. the DB interface stands for a component performing I/O and so should accept a `context` and return an error in all methods; I haven't bothered because of the [My Non-Goals](#my-non-goals)
1. no tests (unless for exploratory reasons)
1. just-enough CSS

   Aesthetic is not a goal here but we also don't want our eyes to bleed; so I just styled it with [PicoCSS](https://picocss.com) with default settings… and semantic HTML is all I need write, sweet.

1. tidy-up templates setup later

   I'm not happy with how I'm parsing and loading templates but I haven't found way that seems idiomatic, logical, and is optimal (parse each template only once and then compose the parsed-trees)

1. git commit messages are a single line and only meant as short-term reminders

## Design Principles

- application code is grouped by entity

  all code that represents that entity lives within that entity folder (incl. http handlers, html templates, etc.)

- template rendering API…

  - template files are embedded `embed.FS` in the application and collected in dedicated package
  - each package of templates exposes one `Write_TemplateName_(io.Writer, TemplateParams)` functions per templates. The function accepts a struct consisting of all the parameters required or supported by the template.
  - a `templates` package sits at the root of the project and contains the common HTML layout code

- let the handler require which kinds of URL it will support (eg. listing, viewing, editing) and what parameters it will expect; let the HTTP server setup code decide the specific URL to use

## Idiomatic Go conventions I've broken

- I often use the name `me` or `my` for the method receiver…
- I have sometimes used else-blocks for the happy path where it's often advised to keep the happy path at "indentation level zero".
