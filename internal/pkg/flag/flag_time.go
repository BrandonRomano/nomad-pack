// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package flag

import (
	"os"
	"strings"
	"time"

	"github.com/posener/complete"
)

// -- DurationVar and durationValue
type DurationVar struct {
	Name       string
	Aliases    []string
	Usage      string
	Default    time.Duration
	Hidden     bool
	EnvVar     string
	Target     *time.Duration
	Completion complete.Predictor
}

type DurationVarP struct {
	*DurationVar
	Shorthand string
}

func (f *Set) DurationVar(i *DurationVar) {
	f.DurationVarP(&DurationVarP{
		DurationVar: i,
		Shorthand:   "",
	})
}

func (f *Set) DurationVarP(i *DurationVarP) {
	initial := i.Default
	if v, exist := os.LookupEnv(i.EnvVar); exist {
		if d, err := time.ParseDuration(appendDurationSuffix(v)); err == nil {
			initial = d
		}
	}

	def := ""
	if i.Default != 0 {
		def = i.Default.String()
	}

	f.VarFlagP(&VarFlagP{
		VarFlag: &VarFlag{
			Name:       i.Name,
			Aliases:    i.Aliases,
			Usage:      i.Usage,
			Default:    def,
			EnvVar:     i.EnvVar,
			Value:      newDurationValue(initial, i.Target, i.Hidden),
			Completion: i.Completion,
		},
		Shorthand: i.Shorthand,
	})
}

type durationValue struct {
	hidden bool
	target *time.Duration
}

func newDurationValue(def time.Duration, target *time.Duration, hidden bool) *durationValue {
	*target = def
	return &durationValue{
		hidden: hidden,
		target: target,
	}
}

func (d *durationValue) Set(s string) error {
	v, err := time.ParseDuration(appendDurationSuffix(s))
	if err != nil {
		return err
	}
	*d.target = v
	return nil
}

func (d *durationValue) Get() interface{} { return *d.target }
func (d *durationValue) String() string   { return d.target.String() }
func (d *durationValue) Example() string  { return "duration" }
func (d *durationValue) Hidden() bool     { return d.hidden }
func (d *durationValue) Type() string     { return "duration" }

// appendDurationSuffix is used as a backwards-compat tool for assuming users
// meant "seconds" when they do not provide a suffixed duration value.
func appendDurationSuffix(s string) string {
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "m") || strings.HasSuffix(s, "h") {
		return s
	}
	return s + "s"
}
