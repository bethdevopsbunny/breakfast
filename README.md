# Breakfast - Hash Manager
the best café for simple password hash lookup<br>
dynamically updated store of popular passwords hashed with popular hashes 



getting some serious clashing on names in this;<br>
for reference;<br>
<br>
wordlist - a text document containing a number of comman passwords <br>
hashes - a set of hashes used on wordlists to give hashed version<br>
storehash - a SHA1 hash of data to help store and retrive data in this app<br>
&nbsp;&nbsp;&nbsp;&nbsp;- full - the whole store object hashed - not sure if im usign this yet<br>
&nbsp;&nbsp;&nbsp;&nbsp;- source - type, owner, repo hashed - used to see if we have a new source in config<br>
&nbsp;&nbsp;&nbsp;&nbsp;- includedwordlists - just included files slice hashed - to check if config has new wordlists <br>
