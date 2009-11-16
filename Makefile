include $(GOROOT)/src/Make.$(GOARCH)

.SUFFIXES: .go .$O

OBJS=irc.$O\
	 main.$O

gobot: $(OBJS)
	$(LD) -o gobot $(OBJS)

clean:
	rm -f $(OBJS) gobot

.go.$O:
	$(GC) $<

