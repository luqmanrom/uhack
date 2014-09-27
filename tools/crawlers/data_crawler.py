import pull_target
x = ["socks", "pants", "shirt", "shoes", "toy", "toys", "blanket", "pillow"]
for i in x:
	pull_target.pull_info("walmart", i)
	pull_target.pull_info("target", i)
	
