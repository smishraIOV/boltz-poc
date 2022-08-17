function waitForLndToSync () {
  while true; do
    if lncli --network=regtest getinfo 2>&1; then
      break
    fi
    sleep 1
  done

  sleep 5
}


while true; do if lncli --network=regtest getinfo 2>&1; then; break; fi; sleep 1; done;
      break
    fi
    sleep 1
  done

  sleep 5



while true; do if true; then; break; fi; sleep 1; done;
