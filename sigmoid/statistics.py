#!/usr/bin/python3

import subprocess
import time

def du(path):
    return subprocess.check_output(['du','-sb', path]).split()[0].decode('utf-8')

if __name__ == "__main__":
    f = open('data.csv', 'w')
    while True:
        string = du('chaindata') + '\t' + du('statedata') +'\n'
        print(string)
        f.write(string)    
        f.flush()
        time.sleep(3)
    f.close()    


