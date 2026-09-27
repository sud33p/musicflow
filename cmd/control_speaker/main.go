package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/mafredri/goodspeaker"
	"github.com/sud33p/musicflow"
	"github.com/sud33p/musicflow/api"
)

var (
	key = "4efgvbn m546Uy7kolKrftgbn =-0u&~"
	iv  = "54eRty@hkL,;/y9U"
)

// PlayCmdRequest is a guess at the structure for PLAY_CMD
type PlayCmdRequest struct {
	// Cmd int `json:"cmd"` // Potential payload?
}

func (PlayCmdRequest) Message() string { return "PLAY_CMD" }

func main() {
	host := flag.String("addr", "192.168.1.156", "Host address or IP of the speaker")
	vol := flag.Int("vol", -1, "Volume level to set (0-100). If -1, volume is not changed.")
	mode := flag.String("mode", "", "Function mode: portable or bluetooth. If empty, mode is not changed.")
	
	flag.Parse()

	if *host == "" {
		log.Fatal("Speaker address must be provided via -addr")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Connecting to %s...", *host)

	var gsOpt []goodspeaker.Option
	if key != "" && iv != "" {
		aes, err := goodspeaker.WithAES([]byte(key), []byte(iv))
		if err != nil {
			log.Fatalf("AES setup failed: %v", err)
		}
		gsOpt = append(gsOpt, aes)
	}
	opt := []musicflow.DialOption{
		musicflow.WithGoodspeakerOption(gsOpt...),
		musicflow.WithLogger(log.New(os.Stderr, "[musicflow] ", log.Flags())),
	}

	c, err := musicflow.Dial(ctx, *host+":9741", opt...)
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}
	defer c.Close()

	// 1. Change Mode
	if *mode != "" {
		var f api.Function
		switch *mode {
		case "bluetooth":
			f = api.FunctionBluetooth
		case "portable":
			f = api.FunctionPortable
		default:
			log.Printf("Unknown mode %q, skipping mode change.", *mode)
		}

		if f != 0 || *mode == "wifi" { // Assuming wifi is 0, checking valid enum mapping if needed or just trusting switch
			if *mode == "portable" || *mode == "bluetooth" {
				log.Printf("Setting mode to %s...", f)
				err = c.Function(ctx, f)
				if err != nil {
					log.Printf("Failed to set mode: %v", err)
				} else {
					log.Printf("Mode set to %s successfully.", f)
				}
			}
		}
	} else {
		log.Println("Mode parameter missing, skipping mode change.")
	}

	// 2. Set Volume
	if *vol != -1 {
		log.Printf("Setting volume to %d...", *vol)
		err = c.Volume(ctx, *vol, 0)
		if err != nil {
			log.Printf("Failed to set volume: %v", err)
		} else {
			log.Println("Volume set successfully.")
		}
	} else {
		log.Println("Volume parameter missing, skipping volume change.")
	}

	// 3. Play/Pause
	// log.Printf("Attempting to send play command (1s timeout)...")
	
	// playCtx, playCancel := context.WithTimeout(ctx, 1*time.Second)
	// err = c.Send(playCtx, musicflow.Request{Message: "PLAY_CMD"}, nil)
	// playCancel()
	// if err != nil {
	// 	if err == context.DeadlineExceeded {
	// 		log.Println("PLAY_CMD sent (no confirmation received)")
	// 	} else {
	// 		log.Printf("Failed to send PLAY_CMD: %v", err)
	// 	}
	// } else {
	// 	log.Println("PLAY_CMD sent.")
	// }

	log.Println("Done.")
	os.Exit(0)
}
