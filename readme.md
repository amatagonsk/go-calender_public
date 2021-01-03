[toc]

# at start

## change hook dir
run
`git config core.hooksPath ./.githooks`

if not using windows, run (idk work)
`chmod a+x ./.githooks/pre-commit`

## build for windows
`fyne package -os windows -release -appID 0.0.1`
