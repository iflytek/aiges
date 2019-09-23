package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Options struct {
	cmd       string
	args      []string
	zkHosts   string
	xmlFile   string
	znodePath string
	depth     int
	force     bool
	debug     bool
}

func parseCmdLine() (*Options, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:

  %s [options] <command> [args]

Command:

  import    Imports the zookeeper tree from XML file. 
            Must be specified with -zookeeper AND -xmlfile options. 
            Optionally takes -path for importing subtree

  export    Exports the zookeeper tree to XML file. 
            Must be specified with -zookeeper option. 
            Optionally takes -path for exporting subtree

  diff      Diff the zookeeper tree based on XML data. 
            Must be specified with -zookeeper OR -xmlfile options. 
            Optionally takes -path for subtree diff

  dump      Dumps the entire ZK (sub)tree to standard output. 
            Must be specified with --zookeeper OR --xmlfile options. 
            Optionally takes --path and --depth for dumping subtree.

  sync      Sync the zookeeper tree with stdin/stdout stream
            Must be specified with -zookeeper option. 
            Optionally takes -path for exporting subtree

Options:

`, os.Args[0])

		flag.PrintDefaults()
	}

	var opts Options

	flag.StringVar(&opts.zkHosts, "zookeeper", "localhost:2181", "specifies information to connect to zookeeper.")
	flag.StringVar(&opts.xmlFile, "xmlfile", "", "Zookeeper tree-data XML file.")
	flag.StringVar(&opts.znodePath, "path", "/", "Path to the zookeeper subtree rootnode.")
	flag.IntVar(&opts.depth, "depth", -1, "Depth of the ZK tree to be dumped (ignored for XML dump).")
	flag.BoolVar(&opts.force, "force", false, "Forces cleanup before import; also used for forceful update.")
	flag.BoolVar(&opts.debug, "debug", false, "Enable debug mode")

	flag.Parse()

	if flag.NArg() == 0 {
		return nil, errors.New("missing command")
	}

	cmd := flag.Arg(0)

	switch cmd {
	case "import", "diff":
		if len(opts.zkHosts) == 0 || len(opts.xmlFile) == 0 {
			return nil, errors.New("missing params")
		}

	case "export", "dump", "sync":
		if len(opts.zkHosts) == 0 {
			return nil, errors.New("missing params")
		}

	default:
		return nil, fmt.Errorf("unknown command: %s", cmd)
	}

	opts.cmd = cmd
	opts.args = flag.Args()[1:]

	return &opts, nil
}

func main() {
	if opts, err := parseCmdLine(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err.Error())

		flag.Usage()

		os.Exit(-1)
	} else {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		if opts.debug {
			log.SetOutput(os.Stderr)
		} else if f, err := os.OpenFile(os.Args[0]+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
			log.SetOutput(os.Stderr)
			log.Fatalf("fail to open log file, %s", err)
		} else {
			log.SetOutput(f)

			defer f.Sync()
		}

		switch opts.cmd {
		case "import":
			if liveTree, err := NewZkTree(strings.Split(opts.zkHosts, ";"), opts.znodePath); err != nil {
				log.Fatalf("fail to connect %s, %s", opts.zkHosts, err)
			} else if loadedTree, err := LoadZkTree(opts.xmlFile); err != nil {
				log.Fatalf("fail to load from %s, %s", opts.xmlFile, err)
			} else if err := liveTree.Merge(loadedTree, opts.force); err != nil {
				log.Fatalf("fail to write to %s, %s", opts.znodePath, err)
			} else {
				log.Printf("import %d nodes from XML file `%s` successful!", loadedTree.Root().Len(), opts.xmlFile)

				if xml, err := loadedTree.Xml(); err == nil {
					log.Printf("%s", string(xml))
				}
			}

		case "export":
			if liveTree, err := NewZkTree(strings.Split(opts.zkHosts, ";"), opts.znodePath); err != nil {
				log.Fatalf("fail to connect %s, %s", opts.zkHosts, err)
			} else if xml, err := liveTree.Xml(); err != nil {
				log.Fatalf("fail to dump XML from %s, %s", opts.znodePath, err)
			} else if len(opts.xmlFile) == 0 {
				os.Stdout.Write(xml)
			} else if err := ioutil.WriteFile(opts.xmlFile, xml, 0644); err != nil {
				log.Fatalf("fail to write XML file `%s`, %s", opts.xmlFile, err)
			} else {
				log.Printf("wrote %d bytes to XML file `%s`", len(xml), opts.xmlFile)
			}

		case "diff":
			if liveTree, err := NewZkTree(strings.Split(opts.zkHosts, ";"), opts.znodePath); err != nil {
				log.Fatalf("fail to connect %s, %s", opts.zkHosts, err)
			} else if loadedTree, err := LoadZkTree(opts.xmlFile); err != nil {
				log.Fatalf("fail to load from %s, %s", opts.xmlFile, err)
			} else if err := liveTree.Diff(loadedTree); err != nil {
				log.Fatalf("fail to diff tree at %s, %s", opts.znodePath, err)
			}

		case "dump":
			var tree ZkTree

			if len(opts.zkHosts) > 0 {
				if liveTree, err := NewZkTree(strings.Split(opts.zkHosts, ";"), opts.znodePath); err != nil {
					log.Fatalf("fail to connect %s, %s", opts.zkHosts, err)
				} else {
					tree = liveTree
				}
			} else if len(opts.xmlFile) > 0 {
				if loadedTree, err := LoadZkTree(opts.xmlFile); err != nil {
					log.Fatalf("fail to load from %s, %s", opts.xmlFile, err)
				} else {
					tree = loadedTree
				}
			}

			if out, err := tree.Dump(opts.depth); err != nil {
				log.Fatalf("fail to dump tree, %s", err)
			} else {
				os.Stdout.WriteString(out)
			}

		case "sync":
			if liveTree, err := NewZkTree(strings.Split(opts.zkHosts, ";"), opts.znodePath); err != nil {
				log.Fatalf("fail to connect %s, %s", opts.zkHosts, err)
			} else if err := liveTree.Sync(os.Stdin, os.Stdout); err != nil {
				log.Fatalf("fail to sync with input #%v and output #%v, %s", os.Stdin.Fd(), os.Stdout.Fd(), err)
			}
		}
	}
}
