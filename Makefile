.PHONY: install
install:
	mkdir -p ~/.config/synth
	touch ~/.config/synth/config.yaml
	@echo 'root-path: '$${HOME}'/synth' > ~/.config/synth/config.yaml
	@echo 'sample-rate: 44100' >> ~/.config/synth/config.yaml
	go install