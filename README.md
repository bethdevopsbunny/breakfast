# Breakfast - Hash Manager
the best café for simple password hash lookup<br>
dynamically updated store of popular passwords hashed with popular hashes 


i use github issues to remember ideas more then anything else 


getting some serious clashing on names in this;
for reference;

wordlist - a text document containing a number of comman passwords
hashes - a set of hashes used on wordlists to give hashed version
storehash - a SHA1 hash of data to help store and retrive data in this app
    - full - the whole store object hashed - not sure if im usign this yet
    - source - type, owner, repo hashed - used to see if we have a new source in config
    - includedwordlists - just included files slice hashed - to check if config has new wordlists 