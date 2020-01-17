/*
 * Warp (C) 2019-2020 MinIO, Inc.
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
 */

package cli

import (
	"github.com/minio/cli"
	"github.com/minio/mc/pkg/probe"
	"github.com/minio/minio-go/v6"
	"github.com/minio/warp/pkg/bench"
)

var (
	mixedFlags = []cli.Flag{
		cli.IntFlag{
			Name:  "objects",
			Value: 2500,
			Usage: "Number of objects to upload.",
		},
		cli.StringFlag{
			Name:  "obj.size",
			Value: "10MB",
			Usage: "Size of each generated object. Can be a number or 10KB/MB/GB. All sizes are base 2 binary.",
		},
		cli.Float64Flag{
			Name:  "get-distrib",
			Usage: "The amount of GET operations.",
			Value: 45,
		},
		cli.Float64Flag{
			Name:  "stat-distrib",
			Usage: "The amount of STAT operations.",
			Value: 30,
		},
		cli.Float64Flag{
			Name:  "put-distrib",
			Usage: "The amount of PUT operations.",
			Value: 15,
		},
		cli.Float64Flag{
			Name:  "delete-distrib",
			Usage: "The amount of DELETE operations. Must be at least the same as PUT.",
			Value: 10,
		},
	}
)

var mixedCmd = cli.Command{
	Name:   "mixed",
	Usage:  "benchmark mixed objects",
	Action: mainMixed,
	Before: setGlobalsFromContext,
	Flags:  combineFlags(globalFlags, ioFlags, mixedFlags, genFlags, benchFlags, analyzeFlags),
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}} [FLAGS]

FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}

EXAMPLES:
...
 `,
}

// mainMixed is the entry point for mixed command.
func mainMixed(ctx *cli.Context) error {
	checkMixedSyntax(ctx)
	src := newGenSource(ctx)
	sse := newSSE(ctx)
	dist := bench.MixedDistribution{
		Distribution: map[string]float64{
			"GET":    ctx.Float64("get-distrib"),
			"STAT":   ctx.Float64("stat-distrib"),
			"PUT":    ctx.Float64("put-distrib"),
			"DELETE": ctx.Float64("delete-distrib"),
		},
	}
	err := dist.Generate(ctx.Int("objects") * 2)
	fatalIf(probe.NewError(err), "Invalid distribution")
	b := bench.Mixed{
		Common: bench.Common{
			Client:      newClient(ctx),
			Concurrency: ctx.Int("concurrent"),
			Source:      src,
			Bucket:      ctx.String("bucket"),
			Location:    "",
			PutOpts: minio.PutObjectOptions{
				ServerSideEncryption: sse,
			},
		},
		CreateObjects: ctx.Int("objects"),
		GetOpts:       minio.GetObjectOptions{ServerSideEncryption: sse},
		StatOpts: minio.StatObjectOptions{
			GetObjectOptions: minio.GetObjectOptions{
				ServerSideEncryption: sse,
			},
		},
		Dist: &dist,
	}
	return runBench(ctx, &b)
}

func checkMixedSyntax(ctx *cli.Context) {
	checkAnalyze(ctx)
	checkBenchmark(ctx)
}