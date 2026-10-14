/**
 * Copyright 2026 AiT authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/

package keepalive

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/os-artificer/AiT/internal/pkg/errhint"
	pkghelpdoc "github.com/os-artificer/AiT/internal/pkg/helpdoc"
	"github.com/os-artificer/AiT/internal/pkg/toolname"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var commandOrder = []string{"add", "reload", "remove", "serve", "start", "status", "stop"}

var commandShort = map[string]string{
	"add":    "Add or replace one unit in keepalive.yaml",
	"reload": "Reread keepalive.yaml and reconcile running units",
	"remove": "Remove one unit from keepalive.yaml",
	"serve":  "Listen for GET /ping and supervise units",
	"start":  "Start one unit",
	"status": "Show unit state and current samples",
	"stop":   "Stop one unit",
}

const (
	helpFlagNameWidth = 18
	cobraHelpFlag     = "help"
)

func registered(name string) bool {
	_, ok := commandShort[name]
	return ok
}

func newKeepCommand(name string) *cobra.Command {
	switch name {
	case "serve":
		return newServeCmd()
	case "status":
		return newStatusCmd()
	case "start":
		return newNameCmd("start", commandShort["start"], NotifyStart)
	case "stop":
		return newNameCmd("stop", commandShort["stop"], NotifyStop)
	case "reload":
		return newReloadCmd()
	case "add":
		return newAddCmd()
	case "remove":
		return newRemoveCmd()
	default:
		return &cobra.Command{Use: name}
	}
}

func newServeCmd() *cobra.Command {
	var listen string
	cmd := &cobra.Command{
		Use:   "serve",
		Short: commandShort["serve"],
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := appRoot()
			if err != nil {
				return err
			}
			parent := cmd.Context()
			if parent == nil {
				parent = context.Background()
			}
			ctx, stop := signal.NotifyContext(parent, os.Interrupt)
			defer stop()
			return Serve(ctx, root, listen)
		},
	}
	cmd.Flags().StringVar(&listen, "listen", "", "HTTP listen address (default "+DefaultListen+")")
	bindHelp(cmd, toolname.KeepCLI("serve --listen "+DefaultListen))
	return cmd
}

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: commandShort["status"],
		RunE: func(cmd *cobra.Command, _ []string) error {
			root, err := appRoot()
			if err != nil {
				return err
			}
			text, err := QueryStatus(root)
			if err != nil {
				return daemonHint(err)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), text)
			return err
		},
	}
	bindHelp(cmd, toolname.KeepCLI("status"))
	return cmd
}

func newReloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: commandShort["reload"],
		RunE: func(*cobra.Command, []string) error {
			root, err := appRoot()
			if err != nil {
				return err
			}
			return daemonHint(NotifyReload(root))
		},
	}
	bindHelp(cmd, toolname.KeepCLI("reload"))
	return cmd
}

func newNameCmd(use, short string, fn func(root, name string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   use + " NAME",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := appRoot()
			if err != nil {
				return err
			}
			return daemonHint(fn(root, args[0]))
		},
	}
	bindHelp(cmd, toolname.KeepCLI(use+" demo"))
	return cmd
}

func newRemoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove NAME",
		Short: commandShort["remove"],
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := appRoot()
			if err != nil {
				return err
			}
			return removeUnit(root, args[0])
		},
	}
	bindHelp(cmd, toolname.KeepCLI("remove demo"))
	return cmd
}

func newAddCmd() *cobra.Command {
	var name, cwd string
	var cpu, diskPct, ioPct float64
	var mem, diskBytes, readBPS, writeBPS uint64
	cmd := &cobra.Command{
		Use: "add --name NAME [--cwd DIR] [--cpu P] [--mem BYTES] " +
			"[--disk-bytes N] [--disk-percent P] [--read-bps N] [--write-bps N] [--io-percent P] -- COMMAND...",
		Short: commandShort["add"],
		RunE: func(_ *cobra.Command, args []string) error {
			root, err := appRoot()
			if err != nil {
				return err
			}
			unit := Unit{
				Name: name, Command: args, Cwd: cwd,
				Limits: Limits{
					CPUPercent: cpu, MemoryBytes: mem,
					DiskSpaceBytes: diskBytes, DiskSpacePercent: diskPct,
					DiskReadBPS: readBPS, DiskWriteBPS: writeBPS, DiskIOPercent: ioPct,
				},
			}
			return addUnit(root, unit)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "unit name")
	cmd.Flags().StringVar(&cwd, "cwd", "", "working directory counted for disk limits")
	cmd.Flags().Float64Var(&cpu, "cpu", 0, "CPU limit as percent of one core")
	cmd.Flags().Uint64Var(&mem, "mem", 0, "RSS limit in bytes")
	cmd.Flags().Uint64Var(&diskBytes, "disk-bytes", 0, "cwd tree byte limit")
	cmd.Flags().Float64Var(&diskPct, "disk-percent", 0, "filesystem used percent limit")
	cmd.Flags().Uint64Var(&readBPS, "read-bps", 0, "process read bytes per second")
	cmd.Flags().Uint64Var(&writeBPS, "write-bps", 0, "process write bytes per second")
	cmd.Flags().Float64Var(&ioPct, "io-percent", 0, "device IO busy percent")
	bindHelp(cmd, toolname.KeepCLI("add --name demo --cpu 20 -- sleep 30"))
	return cmd
}

func bindHelp(cmd *cobra.Command, examples ...string) {
	name := cmd.Name()
	usage := cmd.Use
	flags := formatCmdFlags(cmd)
	cmd.SilenceUsage = true
	cmd.Long = fmt.Sprintf("Run '%s --help' for usage and examples.", toolname.KeepCLI(name))
	pkghelpdoc.Bind(cmd, func(w io.Writer) error {
		return writeCommandHelp(w, name, usage, flags, examples...)
	})
}

func formatCmdFlags(cmd *cobra.Command) string {
	var b strings.Builder
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == cobraHelpFlag {
			return
		}
		label := "--" + f.Name
		if f.Shorthand != "" {
			label = "-" + f.Shorthand + ", " + label
		}
		fmt.Fprintf(&b, "  %-*s %s\n", helpFlagNameWidth, label, f.Usage)
	})
	return b.String()
}

func addUnit(root string, unit Unit) error {
	cfg, _, _ := Layout(root)
	doc, err := Load(cfg)
	if err != nil {
		return err
	}
	doc, err = Replace(doc, unit)
	if err != nil {
		return unitHint(err)
	}
	if err := Save(cfg, doc); err != nil {
		return err
	}
	if err := NotifyStart(root, unit.Name); err != nil && !errors.Is(err, ErrDaemonDown) {
		return err
	}
	return nil
}

func removeUnit(root, name string) error {
	cfg, _, _ := Layout(root)
	doc, err := Load(cfg)
	if err != nil {
		return err
	}
	doc, err = Remove(doc, name)
	if err != nil {
		return unitHint(err)
	}
	if err := Save(cfg, doc); err != nil {
		return err
	}
	if err := NotifyStop(root, name); err != nil && !errors.Is(err, ErrDaemonDown) {
		return err
	}
	return nil
}

func daemonHint(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrDaemonDown) {
		return errhint.WithTry(err, "keepalive serve is not running.", toolname.KeepCLI("serve"))
	}
	return err
}

func unitHint(err error) error {
	switch {
	case errors.Is(err, ErrBadName):
		return errhint.WithTry(err, "Unit names are letters, digits, '-' and '_'.", toolname.KeepCLI("add --name demo -- sleep 30"))
	case errors.Is(err, ErrNoCommand):
		return errhint.WithTry(err, "Pass a command after --.", toolname.KeepCLI("add --name demo -- sleep 30"))
	case errors.Is(err, ErrBadLimit):
		return errhint.WithTry(err, "disk-percent and io-percent must be between 0 and 100.", toolname.KeepCLI("add --name demo --disk-percent 90 -- sleep 30"))
	case errors.Is(err, ErrUnitMissing):
		return errhint.WithTry(err, "That unit is not in keepalive.yaml.", toolname.KeepCLI("status"), toolname.KeepCLI("add --name demo -- sleep 30"))
	case errors.Is(err, ErrUnitExists):
		return errhint.WithTry(err, "A unit with that name already exists.", toolname.KeepCLI("remove NAME"), toolname.KeepCLI("add --name demo -- sleep 30"))
	default:
		return err
	}
}
