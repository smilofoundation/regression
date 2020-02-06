// Copyright 2020 smilofoundation/regression Authors
// Copyright 2019 smilofoundation/regression Authors
// Copyright 2017 AMIS Technologies
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package compose

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"go-smilo/src/blockchain/regression/src/docker/service"
)

type Compose interface {
	String() string
}

type sport struct {
	IPPrefix string
	EthStats *service.EthStats
	Services []*service.Fullnode
}

func New(ipPrefix string, number int, secret string, nodeKeys []string,
	genesis string, staticNodes string, smilo bool) Compose {
	ist := &sport{
		IPPrefix: ipPrefix,
		EthStats: service.NewEthStats(fmt.Sprintf("%v.9", ipPrefix), secret),
	}
	ist.init(number, nodeKeys, genesis, staticNodes)
	if smilo {
		return newSmilo(ist, number)
	}
	return ist
}

func (ist *sport) init(number int, nodeKeys []string, genesis string, staticNodes string) {
	for i := 0; i < number; i++ {
		s := service.NewFullnode(i,
			genesis,
			nodeKeys[i],
			"",
			30303+i,
			8545+i,
			ist.EthStats.Host(),
			// from subnet ip 10
			fmt.Sprintf("%v.%v", ist.IPPrefix, i+10),
		)

		staticNodes = strings.Replace(staticNodes, "0.0.0.0", s.IP, 1)
		ist.Services = append(ist.Services, s)
	}

	// update static nodes
	for i := range ist.Services {
		ist.Services[i].StaticNodes = staticNodes
	}
}

func (ist sport) String() string {
	tmpl, err := template.New("sport").Parse(sportTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, ist)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

var sportTemplate = `version: '3'
services:
  {{ .EthStats }}
  {{- range .Services }}
  {{ . }}
  {{- end }}
networks:
  app_net:
    driver: bridge
    ipam:
      driver: default
      config:
      - subnet: {{ .IPPrefix }}.0/24`
