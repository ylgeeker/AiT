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

package imgmeta

import (
	"time"

	"github.com/dsoprea/go-exif/v3"
)

func setDateTimeTag(ib *exif.IfdBuilder, name string, value time.Time) error {
	return ib.SetStandardWithName(name, value)
}

func applyPatchExifIFD(rootIb *exif.IfdBuilder, p PatchFields) error {
	needExif := p.LensMake != nil || p.LensModel != nil || p.SerialNumber != nil || p.DateTime != nil
	if !needExif {
		return nil
	}

	exifIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPathExif)
	if err != nil {
		return err
	}
	if err := setASCII(exifIb, "LensMake", p.LensMake); err != nil {
		return err
	}
	if err := setASCII(exifIb, "LensModel", p.LensModel); err != nil {
		return err
	}
	if err := setASCII(exifIb, "BodySerialNumber", p.SerialNumber); err != nil {
		return err
	}
	if p.DateTime == nil {
		return nil
	}

	return setDateTimeTag(exifIb, "DateTimeOriginal", *p.DateTime)
}

func applyPatchGPS(rootIb *exif.IfdBuilder, p PatchFields) error {
	if !p.HasGPS() && p.Altitude == nil {
		return nil
	}

	gpsIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPathGPS)
	if err != nil {
		return err
	}
	if p.HasGPS() {
		if err := setGPSCoords(gpsIb, *p.Latitude, *p.Longitude); err != nil {
			return err
		}
	}
	if p.Altitude == nil {
		return nil
	}

	return setGPSAltitude(gpsIb, *p.Altitude)
}

func clearGPSIFD(gpsIb *exif.IfdBuilder, mask ClearMask) {
	if mask.GPS {
		for _, name := range []string{"GPSLatitude", "GPSLatitudeRef", "GPSLongitude", "GPSLongitudeRef"} {
			_ = deleteTagByName(gpsIb, name)
		}
	}
	if !mask.Altitude {
		return
	}
	for _, name := range []string{"GPSAltitude", "GPSAltitudeRef"} {
		_ = deleteTagByName(gpsIb, name)
	}
}
