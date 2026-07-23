# templates

Rendering of HTML templates from static or embedded filesystems with optional caching of parsed templates, and buffering of template execution to prevent the rendering of partial output.

These templates invert the standard golang template organisation whereby every template includes the contents of header and footer fragments, and instead use a more [Jinja-like](https://jinja.palletsprojects.com/en/stable/templates/#template-inheritance) setup where there is a common layout template that includes the body template provided by each page to be rendered. For an example of usage please review the tests for this package.
