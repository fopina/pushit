package cmd

type CLIOptions struct {
	Profile string
	Output  bool
	Stream  bool
	Tail    int
	WebBind string
	Config  Config
}
