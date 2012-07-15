opts = Variables( 'options.conf', ARGUMENTS )
opts.Add("DESTDIR", 'Set the root directory to install into ( /path/to/DESTDIR )', "")

env = Environment(ENV = {'GOROOT': '/usr/lib/go'}, TOOLS=['default', 'go'],
		  options = opts)

exportedservice = env.Go('exportedservice', ["exportedport.go",
					     "exportedhttpservice.go"])
pack = env.GoPack('exportedservice', exportedservice)

env.Install(env['DESTDIR'] + env['GO_PKGROOT'] +
	    '/ancientsolutions.com/doozer', pack)
env.Alias('install', [env['DESTDIR'] + env['GO_PKGROOT'] +
	  '/ancientsolutions.com/doozer'])

opts.Save('options.conf', env)
