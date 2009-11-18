include $(GOROOT)/src/Make.$(GOARCH)

.SUFFIXES: .go .$O

OBJS=irc.$O\
	 irc_bot.$O\
	 irc_client.$O\
	 main.$O

gobot: $(OBJS)
	$(LD) -o gobot $(OBJS)

clean:
	rm -f $(OBJS) gobot

.go.$O:
	$(GC) $<

