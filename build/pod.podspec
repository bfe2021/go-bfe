Pod::Spec.new do |spec|
  spec.name         = 'Gbfe'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/bfe2021/go-bfe'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Bfedu Client'
  spec.source       = { :git => 'https://github.com/bfe2021/go-bfe.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gbfe.framework'

	spec.prepare_command = <<-CMD
    curl https://gbfestore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gbfe.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
