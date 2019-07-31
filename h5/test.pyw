#get pid of self
#if at correct time start cmd
#run server in half hour
#kill server by pid

import subprocess
import time
import os

if __name__ == "__main__":
	p = os.path.join(os.environ["USERPROFILE"], "Downloads")
	os.chdir(p)
	si = subprocess.STARTUPINFO()
	si.dwFlags = subprocess.STARTF_USESHOWWINDOW
	#si.wShowWindow = subprocess.SW_HIDE
	cs = subprocess.Popen("python -m http.server 9001", stdout=subprocess.DEVNULL, startupinfo=si)
	#cs = subprocess.Popen("python -m http.server 9001")
	time.sleep(10)
	cs.kill()