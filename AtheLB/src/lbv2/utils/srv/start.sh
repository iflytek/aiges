#!/usr/bin/env bash

srv.exe -dur 1000 -lbname lbv2 -live 1 -max 100 -min 10 -svc xvc -subsvc vc -addr "1.1.1.1:1111"

