#!/bin/bash
ls -lah | fd .yaml
for i in {1..3}; do
	echo 'hello';
done
