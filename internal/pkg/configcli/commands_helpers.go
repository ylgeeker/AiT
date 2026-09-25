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

package configcli

import (
	"fmt"
	"io"
	"strings"

	pkgconfig "github.com/os-artificer/AiT/internal/pkg/config"
	"github.com/os-artificer/AiT/internal/pkg/creds"

	"github.com/spf13/cobra"
)

func runSet(cmd *cobra.Command, o Options, key, value string) error {
	if err := o.RejectKey(key); err != nil {
		return err
	}
	if o.HandleAPIKeys {
		if handled, err := trySetAPIKey(cmd, key, value); handled {
			return err
		}
	}

	cfg, err := pkgconfig.Load()
	if err != nil {
		return err
	}
	if err := cfg.SetKey(key, value); err != nil {
		return err
	}
	return pkgconfig.Save(cfg)
}

func runGet(cmd *cobra.Command, o Options, key string) error {
	if err := o.RejectKey(key); err != nil {
		return err
	}
	if o.HandleAPIKeys {
		if handled, err := tryGetAPIKeyMasked(key); handled {
			return err
		}
	}

	cfg, err := pkgconfig.Load()
	if err != nil {
		return err
	}
	v, err := cfg.GetKey(key)
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

func trySetAPIKey(cmd *cobra.Command, key, value string) (handled bool, err error) {
	provider, ok := pkgconfig.ParseAPIKeyConfigKey(key)
	if !ok {
		return false, nil
	}
	if err := creds.SetProviderAPIKey(provider, value); err != nil {
		return true, err
	}
	fmt.Fprintln(cmd.OutOrStdout(), pkgconfig.APIKeyStoredMessage(provider))
	return true, nil
}

func tryGetAPIKeyMasked(key string) (handled bool, err error) {
	provider, ok := pkgconfig.ParseAPIKeyConfigKey(key)
	if !ok {
		return false, nil
	}
	v, err := creds.GetProviderAPIKey(provider)
	if err != nil {
		return true, err
	}
	fmt.Println(pkgconfig.MaskSecret(v))
	return true, nil
}

func printEnsureResult(out io.Writer, res pkgconfig.EnsureFileResult, dryRun bool) {
	switch {
	case res.Created && !dryRun:
		fmt.Fprintf(out, "Created %s\n", res.Path)
	case res.Created && dryRun:
		fmt.Fprintf(out, "Would create %s\n", res.Path)
	case res.Updated:
		if dryRun {
			fmt.Fprintf(out, "Would update %s (added: %s)\n", res.Path, strings.Join(res.AddedKeys, ", "))
			return
		}
		if len(res.AddedKeys) > 0 {
			fmt.Fprintf(out, "Updated %s (added: %s)\n", res.Path, strings.Join(res.AddedKeys, ", "))
			return
		}
		fmt.Fprintf(out, "Updated %s\n", res.Path)
	default:
		fmt.Fprintf(out, "Config at %s is up to date\n", res.Path)
	}
}
