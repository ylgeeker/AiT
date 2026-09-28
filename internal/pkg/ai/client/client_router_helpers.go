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

package client

import (
	"fmt"

	"github.com/os-artificer/AiT/internal/pkg/ai/provider"
	"github.com/os-artificer/AiT/internal/pkg/ai/router"
	"github.com/os-artificer/AiT/internal/pkg/config"
	"github.com/os-artificer/AiT/internal/pkg/errhint"
)

func newClientProviderOverrideErr(cfg *config.Config, override ProviderID, slots map[ProviderID]*providerSlot) error {
	if override == "" {
		return nil
	}
	if !config.IsCloudProvider(override) {
		return fmt.Errorf(errProviderOverrideCloudFmt, ProviderDeepseek, ProviderMinimax)
	}
	if !cfg.ProviderEnabled(override) {
		return ProviderDisabledErr(string(override))
	}
	if _, ok := slots[override]; ok {
		return nil
	}
	return errhint.WithTryMsg(fmt.Sprintf(errAPIKeyNotSetFmt, override),
		provider.HintStoreAPIKeyInKeyring,
		provider.ConfigSetAPIKeyCmd(string(override)),
		provider.ConfigGetAPIKeyCmd(string(override)),
	)
}

func newClientDefaultCloudKeyErr(cfg *config.Config, override ProviderID, slots map[ProviderID]*providerSlot) error {
	if override != "" {
		return nil
	}
	cloud := cfg.DefaultProviderID()
	if !hasCloudSlot(slots) {
		return nil
	}
	if _, ok := slots[cloud]; ok {
		return nil
	}
	return fmt.Errorf(errAPIKeyNotSetWithHintFmt, cloud, router.APIKeyHint(cloud))
}
