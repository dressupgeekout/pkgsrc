package main

import (
	"fmt"

	"gopkg.in/check.v1"
)

func (s *Suite) Test_VartypeCheck_AwkCommand(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).AwkCommand)

	vt.Varname("PLIST_AWK")
	vt.Op(opAssignAppend)
	vt.Values(
		"{print $0}",
		"{print $$0}")

	vt.Output(
		"WARN: fname:1: $0 is ambiguous. Use ${0} if you mean a Makefile variable or $$0 if you mean a shell variable.")
}

func (s *Suite) Test_VartypeCheck_BasicRegularExpression(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).BasicRegularExpression)

	vt.Varname("REPLACE_FILES.pl")
	vt.Values(
		".*\\.pl$",
		".*\\.pl$$")

	vt.Output(
		"WARN: fname:1: Pkglint parse error in MkLine.Tokenize at \"$\".")
}

func (s *Suite) Test_VartypeCheck_BuildlinkDepmethod(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).BuildlinkDepmethod)

	vt.Varname("BUILDLINK_DEPMETHOD.libc")
	vt.Op(opAssignDefault)
	vt.Values(
		"full",
		"unknown")

	vt.Output(
		"WARN: fname:2: Invalid dependency method \"unknown\". Valid methods are \"build\" or \"full\".")
}

func (s *Suite) Test_VartypeCheck_Category(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).Category)

	t.SetupFileLines("filesyscategory/Makefile",
		"# empty")
	t.SetupFileLines("wip/Makefile",
		"# empty")

	vt.Varname("CATEGORIES")
	vt.Values(
		"chinese",
		"arabic",
		"filesyscategory",
		"wip")

	vt.Output(
		"ERROR: fname:2: Invalid category \"arabic\".",
		"ERROR: fname:4: Invalid category \"wip\".")
}

func (s *Suite) Test_VartypeCheck_CFlag(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).CFlag)

	vt.Varname("CFLAGS")
	vt.Op(opAssignAppend)
	vt.Values(
		"-Wall",
		"/W3",
		"target:sparc64",
		"-std=c99",
		"-XX:+PrintClassHistogramAfterFullGC",
		"`pkg-config pidgin --cflags`")

	vt.Output(
		"WARN: fname:2: Compiler flag \"/W3\" should start with a hyphen.",
		"WARN: fname:3: Compiler flag \"target:sparc64\" should start with a hyphen.",
		"WARN: fname:5: Unknown compiler flag \"-XX:+PrintClassHistogramAfterFullGC\".")

	vt.Op(opUseMatch)
	vt.Values(
		"anything")

	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_Comment(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).Comment)

	G.Pkg = NewPackage(t.File("category/converter"))
	G.Pkg.EffectivePkgbase = "converter"

	vt.Varname("COMMENT")
	vt.Values(
		"Versatile Programming Language",
		"TODO: Short description of the package",
		"A great package.",
		"some packages need a very very long comment to explain their basic usefulness",
		"\"Quoting the comment is wrong\"",
		"'Quoting the comment is wrong'",
		"Package is a great package",
		"Package is an awesome package",
		"The Big New Package is a great package",
		"Converter converts between measurement units")

	vt.Output(
		"ERROR: fname:2: COMMENT must be set.",
		"WARN: fname:3: COMMENT should not begin with \"A\".",
		"WARN: fname:3: COMMENT should not end with a period.",
		"WARN: fname:4: COMMENT should start with a capital letter.",
		"WARN: fname:4: COMMENT should not be longer than 70 characters.",
		"WARN: fname:5: COMMENT should not be enclosed in quotes.",
		"WARN: fname:6: COMMENT should not be enclosed in quotes.",
		"WARN: fname:7: COMMENT should not contain \"is a\".",
		"WARN: fname:8: COMMENT should not contain \"is an\".",
		"WARN: fname:9: COMMENT should not contain \"is a\".",
		"WARN: fname:10: COMMENT should not start with the package name.")
}

func (s *Suite) Test_VartypeCheck_ConfFiles(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).ConfFiles)

	vt.Varname("CONF_FILES")
	vt.Op(opAssignAppend)
	vt.Values(
		"single/file",
		"share/etc/config ${PKG_SYSCONFDIR}/etc/config",
		"share/etc/config ${PKG_SYSCONFBASE}/etc/config file",
		"share/etc/config ${PREFIX}/etc/config share/etc/config2 ${VARBASE}/config2",
		"share/etc/bootrc /etc/bootrc")

	vt.Output(
		"WARN: fname:1: Values for CONF_FILES should always be pairs of paths.",
		"WARN: fname:3: Values for CONF_FILES should always be pairs of paths.",
		"WARN: fname:5: Found absolute pathname: /etc/bootrc",
		"WARN: fname:5: The destination file \"/etc/bootrc\" should start with a variable reference.")
}

func (s *Suite) Test_VartypeCheck_Dependency(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Dependency)

	vt.Varname("CONFLICTS")
	vt.Op(opAssignAppend)
	vt.Values(
		"Perl",
		"perl5>=5.22",
		"perl5-*",
		"perl5-5.22.*",
		"perl5-[5.10-5.22]*",
		"py-docs",
		"perl5-5.22.*{,nb*}",
		"libkipi>=0.1.5<4.0",
		"gtk2+>=2.16",
		"perl-5.22",
		"perl-5*",
		"gtksourceview-sharp-2.0-[0-9]*",
		"perl-5.22{,nb*}",
		"perl-5.22{,nb[0-9]*}",
		"mbrola-301h{,nb[0-9]*}",
		"mpg123{,-esound,-nas}>=0.59.18",
		"mysql*-{client,server}-[0-9]*",
		"postgresql8[0-35-9]-${module}-[0-9]*",
		"ncurses-${NC_VERS}{,nb*}",
		"{ssh{,6}-[0-9]*,openssh-[0-9]*}",
		"gnome-control-center>=2.20.1{,nb*}")

	vt.Output(
		"WARN: fname:1: Unknown dependency pattern \"Perl\".",
		"WARN: fname:3: Please use \"perl5-[0-9]*\" instead of \"perl5-*\".",
		"WARN: fname:5: Only [0-9]* is allowed in the numeric part of a dependency.",
		"WARN: fname:5: The version pattern \"[5.10-5.22]*\" should not contain a hyphen.",
		"WARN: fname:6: Unknown dependency pattern \"py-docs\".",
		"WARN: fname:10: Please use \"5.22{,nb*}\" instead of \"5.22\" as the version pattern.",
		"WARN: fname:11: Please use \"5.*\" instead of \"5*\" as the version pattern.",
		"WARN: fname:12: The version pattern \"2.0-[0-9]*\" should not contain a hyphen.",
		"WARN: fname:20: The version pattern \"[0-9]*,openssh-[0-9]*}\" should not contain a hyphen.", // XXX
		"WARN: fname:21: Dependency patterns of the form pkgbase>=1.0 don't need the \"{,nb*}\" extension.")
}

func (s *Suite) Test_VartypeCheck_DependencyWithPath(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).DependencyWithPath)

	t.CreateFileLines("x11/alacarte/Makefile")
	t.CreateFileLines("category/package/Makefile")
	t.CreateFileLines("devel/gettext/Makefile")
	t.CreateFileLines("devel/gmake/Makefile")
	G.Pkg = NewPackage(t.File("category/package"))

	vt.Varname("DEPENDS")
	vt.Op(opAssignAppend)
	vt.File(G.Pkg.File("fname.mk"))
	vt.Values(
		"Perl",
		"perl5>=5.22:../perl5",
		"perl5>=5.24:../../lang/perl5",
		"broken0.12.1:../../x11/alacarte",
		"broken[0-9]*:../../x11/alacarte",
		"broken[0-9]*../../x11/alacarte",
		"broken>=:../../x11/alacarte",
		"broken=0:../../x11/alacarte",
		"broken=:../../x11/alacarte",
		"broken-:../../x11/alacarte",
		"broken>:../../x11/alacarte",
		"gtk2+>=2.16:../../x11/alacarte",
		"gettext-[0-9]*:../../devel/gettext",
		"gmake-[0-9]*:../../devel/gmake")

	vt.Output(
		"WARN: ~/category/package/fname.mk:1: Unknown dependency pattern with path \"Perl\".",
		"WARN: ~/category/package/fname.mk:2: Dependencies should have the form \"../../category/package\".",
		"ERROR: ~/category/package/fname.mk:3: \"../../lang/perl5\" does not exist.",
		"ERROR: ~/category/package/fname.mk:3: There is no package in \"lang/perl5\".",
		"WARN: ~/category/package/fname.mk:3: Please use USE_TOOLS+=perl:run instead of this dependency.",
		"WARN: ~/category/package/fname.mk:4: Unknown dependency pattern \"broken0.12.1\".",
		"WARN: ~/category/package/fname.mk:5: Unknown dependency pattern \"broken[0-9]*\".",
		"WARN: ~/category/package/fname.mk:6: Unknown dependency pattern with path \"broken[0-9]*../../x11/alacarte\".",
		"WARN: ~/category/package/fname.mk:7: Unknown dependency pattern \"broken>=\".",
		"WARN: ~/category/package/fname.mk:8: Unknown dependency pattern \"broken=0\".",
		"WARN: ~/category/package/fname.mk:9: Unknown dependency pattern \"broken=\".",
		"WARN: ~/category/package/fname.mk:10: Unknown dependency pattern \"broken-\".",
		"WARN: ~/category/package/fname.mk:11: Unknown dependency pattern \"broken>\".",
		"WARN: ~/category/package/fname.mk:13: Please use USE_TOOLS+=msgfmt instead of this dependency.",
		"WARN: ~/category/package/fname.mk:14: Please use USE_TOOLS+=gmake instead of this dependency.")
}

func (s *Suite) Test_VartypeCheck_DistSuffix(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).DistSuffix)

	vt.Varname("EXTRACT_SUFX")
	vt.Values(
		".tar.gz",
		".tar.bz2")

	vt.Output(
		"NOTE: fname:1: EXTRACT_SUFX is \".tar.gz\" by default, so this definition may be redundant.")
}

func (s *Suite) Test_VartypeCheck_EmulPlatform(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).EmulPlatform)

	vt.Varname("EMUL_PLATFORM")
	vt.Values(
		"linux-i386",
		"nextbsd-8087",
		"${LINUX}")

	vt.Output(
		"WARN: fname:2: \"nextbsd\" is not valid for the operating system part of EMUL_PLATFORM. "+
			"Use one of "+
			"{ bitrig bsdos cygwin darwin dragonfly freebsd haiku hpux "+
			"interix irix linux mirbsd netbsd openbsd osf1 solaris sunos "+
			"} instead.",
		"WARN: fname:2: \"8087\" is not valid for the hardware architecture part of EMUL_PLATFORM. "+
			"Use one of "+
			"{ aarch64 aarch64eb alpha amd64 arc arm arm26 arm32 cobalt coldfire convex dreamcast "+
			"earm earmeb earmhf earmhfeb earmv4 earmv4eb earmv5 earmv5eb earmv6 earmv6eb earmv6hf "+
			"earmv6hfeb earmv7 earmv7eb earmv7hf earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 "+
			"i386 i586 i686 ia64 m68000 m68k m88k mips mips64 mips64eb mips64el mipseb mipsel mipsn32 "+
			"mlrisc ns32k pc532 pmax powerpc powerpc64 rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 "+
			"} instead.",
		"WARN: fname:3: \"${LINUX}\" is not a valid emulation platform.")
}

func (s *Suite) Test_VartypeCheck_Enum(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), enum("jdk1 jdk2 jdk4").checker)

	vt.Varname("JDK")
	vt.Op(opUseMatch)
	vt.Values(
		"*",
		"jdk*",
		"sun-jdk*",
		"${JDKNAME}",
		"[")

	vt.Output(
		"WARN: fname:3: The pattern \"sun-jdk*\" cannot match any of { jdk1 jdk2 jdk4 } for JDK.",
		"WARN: fname:5: Invalid match pattern \"[\".",
		"WARN: fname:5: The pattern \"[\" cannot match any of { jdk1 jdk2 jdk4 } for JDK.")
}

func (s *Suite) Test_VartypeCheck_Enum__use_match(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	mklines := t.NewMkLines("module.mk",
		MkRcsID,
		"",
		".if !empty(MACHINE_ARCH:Mi386) || ${MACHINE_ARCH} == i386",
		".endif",
		".if !empty(PKGSRC_COMPILER:Mclang) || ${PKGSRC_COMPILER} == clang",
		".endif",
		".if ${MACHINE_ARCH:Ni386:Nx86_64:Nsparc64}",
		".endif")

	mklines.Check()

	t.CheckOutputLines(
		"NOTE: module.mk:3: MACHINE_ARCH should be compared using == instead of the :M or :N modifier without wildcards.",
		"WARN: module.mk:5: Use ${PKGSRC_COMPILER:Mclang} instead of the == operator.")
}

func (s *Suite) Test_VartypeCheck_FetchURL(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).FetchURL)

	t.SetupMasterSite("MASTER_SITE_GNU", "http://ftp.gnu.org/pub/gnu/")
	t.SetupMasterSite("MASTER_SITE_GITHUB", "https://github.com/")

	vt.Varname("MASTER_SITES")
	vt.Values(
		"https://github.com/example/project/",
		"http://ftp.gnu.org/pub/gnu/bison", // Missing a slash at the end
		"${MASTER_SITE_GNU:=bison}",
		"${MASTER_SITE_INVALID:=subdir/}")

	vt.Output(
		"WARN: fname:1: Please use ${MASTER_SITE_GITHUB:=example/} "+
			"instead of \"https://github.com/example/project/\" "+
			"and run \""+confMake+" help topic=github\" for further tips.",
		"WARN: fname:2: Please use ${MASTER_SITE_GNU:=bison} instead of \"http://ftp.gnu.org/pub/gnu/bison\".",
		"ERROR: fname:3: The subdirectory in MASTER_SITE_GNU must end with a slash.",
		"ERROR: fname:4: The site MASTER_SITE_INVALID does not exist.")

	// PR 46570, keyword gimp-fix-ca
	vt.Values(
		"https://example.org/download.cgi?fname=fname&sha1=12341234")

	t.CheckOutputEmpty()

	vt.Values(
		"http://example.org/distfiles/",
		"http://example.org/download?fname=distfile;version=1.0",
		"http://example.org/download?fname=<distfile>;version=<version>")

	vt.Output(
		"WARN: fname:8: \"http://example.org/download?fname=<distfile>;version=<version>\" is not a valid URL.")
}

func (s *Suite) Test_VartypeCheck_Filename(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Filename)

	vt.Varname("FNAME")
	vt.Values(
		"Filename with spaces.docx",
		"OS/2-manual.txt")

	vt.Output(
		"WARN: fname:1: \"Filename with spaces.docx\" is not a valid filename.",
		"WARN: fname:2: A filename should not contain a slash.")

	vt.Op(opUseMatch)
	vt.Values(
		"Filename with spaces.docx")

	// There's no guarantee that a file name only contains [A-Za-z0-9.].
	// Therefore there are no useful checks in this situation.
	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_Filemask(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Filemask)

	vt.Varname("FNAME")
	vt.Values(
		"Filemask with spaces.docx",
		"OS/2-manual.txt")

	vt.Output(
		"WARN: fname:1: \"Filemask with spaces.docx\" is not a valid filename mask.",
		"WARN: fname:2: A filename mask should not contain a slash.")

	vt.Op(opUseMatch)
	vt.Values(
		"Filemask with spaces.docx")

	// There's no guarantee that a file name only contains [A-Za-z0-9.].
	// Therefore there are no useful checks in this situation.
	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_FileMode(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).FileMode)

	vt.Varname("HIGHSCORE_PERMS")
	vt.Values(
		"u+rwx",
		"0600",
		"1234",
		"12345",
		"${OTHER_PERMS}")

	vt.Output(
		"WARN: fname:1: Invalid file mode \"u+rwx\".",
		"WARN: fname:4: Invalid file mode \"12345\".")

	vt.Op(opUseMatch)
	vt.Values(
		"u+rwx")

	// There's no guarantee that a file name only contains [A-Za-z0-9.].
	// Therefore there are no useful checks in this situation.
	vt.Output(
		"WARN: fname:11: Invalid file mode \"u+rwx\".")
}

func (s *Suite) Test_VartypeCheck_Homepage(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).Homepage)

	vt.Varname("HOMEPAGE")
	vt.Values(
		"${MASTER_SITES}")

	vt.Output(
		"WARN: fname:1: HOMEPAGE should not be defined in terms of MASTER_SITEs.")

	G.Pkg = NewPackage(t.File("category/package"))
	G.Pkg.vars.Define("MASTER_SITES", t.NewMkLine(G.Pkg.File("Makefile"), 5, "MASTER_SITES=\thttps://cdn.NetBSD.org/pub/pkgsrc/distfiles/"))

	vt.Values(
		"${MASTER_SITES}")

	vt.Output(
		"WARN: fname:2: HOMEPAGE should not be defined in terms of MASTER_SITEs. Use https://cdn.NetBSD.org/pub/pkgsrc/distfiles/ directly.")
}

func (s *Suite) Test_VartypeCheck_Identifier(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).Identifier)

	vt.Varname("SUBST_CLASSES")
	vt.Values(
		"${OTHER_VAR}",
		"identifiers cannot contain spaces",
		"id/cannot/contain/slashes")
	vt.Op(opUseMatch)
	vt.Values(
		"[A-Z]",
		"A*B")

	vt.Output(
		"WARN: fname:2: Invalid identifier \"identifiers cannot contain spaces\".",
		"WARN: fname:3: Invalid identifier \"id/cannot/contain/slashes\".",
		"WARN: fname:11: Invalid identifier pattern \"[A-Z]\" for SUBST_CLASSES.")
}

func (s *Suite) Test_VartypeCheck_Integer(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).Integer)

	vt.Varname("NUMBER")
	vt.Values(
		"${OTHER_VAR}",
		"123",
		"-13",
		"11111111111111111111111111111111111111111111111")

	vt.Output(
		"WARN: fname:1: Invalid integer \"${OTHER_VAR}\".",
		"WARN: fname:3: Invalid integer \"-13\".")
}

func (s *Suite) Test_VartypeCheck_LdFlag(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).LdFlag)

	vt.Varname("LDFLAGS")
	vt.Op(opAssignAppend)
	vt.Values(
		"-lc",
		"-L/usr/lib64",
		"`pkg-config pidgin --ldflags`",
		"-unknown",
		"no-hyphen",
		"-Wl,--rpath,/usr/lib64")
	vt.Op(opUseMatch)
	vt.Values(
		"anything")

	vt.Output(
		"WARN: fname:4: Unknown linker flag \"-unknown\".",
		"WARN: fname:5: Linker flag \"no-hyphen\" should start with a hyphen.",
		"WARN: fname:6: Please use \"${COMPILER_RPATH_FLAG}\" instead of \"-Wl,--rpath\".")
}

func (s *Suite) Test_VartypeCheck_License(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).License)

	vt.Varname("LICENSE")
	vt.Values(
		"gnu-gpl-v2",
		"AND mit")

	vt.Output(
		"WARN: fname:1: License file ~/licenses/gnu-gpl-v2 does not exist.",
		"ERROR: fname:2: Parse error for license condition \"AND mit\".")

	vt.Op(opAssignAppend)
	vt.Values(
		"gnu-gpl-v2",
		"AND mit")

	vt.Output(
		"ERROR: fname:11: Parse error for appended license condition \"gnu-gpl-v2\".",
		"WARN: fname:12: License file ~/licenses/mit does not exist.")
}

func (s *Suite) Test_VartypeCheck_MachineGnuPlatform(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).MachineGnuPlatform)

	vt.Varname("MACHINE_GNU_PLATFORM")
	vt.Op(opUseMatch)
	vt.Values(
		"x86_64-pc-cygwin",
		"Cygwin-*-amd64",
		"x86_64-*",
		"*-*-*-*",
		"${OTHER_VAR}")

	vt.Output(
		"WARN: fname:2: The pattern \"Cygwin\" cannot match any of "+
			"{ aarch64 aarch64_be alpha amd64 arc arm armeb armv4 armv4eb armv6 armv6eb armv7 armv7eb "+
			"cobalt convex dreamcast hpcmips hpcsh hppa hppa64 i386 i486 ia64 m5407 m68010 m68k m88k "+
			"mips mips64 mips64el mipseb mipsel mipsn32 mlrisc ns32k pc532 pmax powerpc powerpc64 "+
			"rs6000 s390 sh shle sparc sparc64 vax x86_64 "+
			"} for the hardware architecture part of MACHINE_GNU_PLATFORM.",
		"WARN: fname:2: The pattern \"amd64\" cannot match any of "+
			"{ bitrig bsdos cygwin darwin dragonfly freebsd haiku hpux interix irix linux mirbsd "+
			"netbsd openbsd osf1 solaris sunos } "+
			"for the operating system part of MACHINE_GNU_PLATFORM.",
		"WARN: fname:4: \"*-*-*-*\" is not a valid platform pattern.")
}

func (s *Suite) Test_VartypeCheck_MailAddress(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).MailAddress)

	vt.Varname("MAINTAINER")
	vt.Values(
		"pkgsrc-users@netbsd.org",
		"tech-pkg@NetBSD.org",
		"packages@NetBSD.org",
		"user1@example.org,user2@example.org")

	vt.Output(
		"WARN: fname:1: Please write \"NetBSD.org\" instead of \"netbsd.org\".",
		"ERROR: fname:2: This mailing list address is obsolete. Use pkgsrc-users@NetBSD.org instead.",
		"ERROR: fname:3: This mailing list address is obsolete. Use pkgsrc-users@NetBSD.org instead.",
		"WARN: fname:4: \"user1@example.org,user2@example.org\" is not a valid mail address.")
}

func (s *Suite) Test_VartypeCheck_Message(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Message)

	vt.Varname("SUBST_MESSAGE.id")
	vt.Values(
		"\"Correct paths\"",
		"Correct paths")

	vt.Output(
		"WARN: fname:1: SUBST_MESSAGE.id should not be quoted.")
}

func (s *Suite) Test_VartypeCheck_Option(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Option)

	G.Pkgsrc.PkgOptions["documented"] = "Option description"
	G.Pkgsrc.PkgOptions["undocumented"] = ""

	vt.Varname("PKG_OPTIONS.pkgbase")
	vt.Values(
		"documented",
		"undocumented",
		"unknown",
		"underscore_is_deprecated",
		"UPPER")

	vt.Output(
		"WARN: fname:3: Unknown option \"unknown\".",
		"WARN: fname:4: Use of the underscore character in option names is deprecated.",
		"ERROR: fname:5: Invalid option name \"UPPER\". "+
			"Option names must start with a lowercase letter and be all-lowercase.")
}

func (s *Suite) Test_VartypeCheck_Pathlist(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Pathlist)

	vt.Varname("PATH")
	vt.Values(
		"/usr/bin:/usr/sbin:.::${LOCALBASE}/bin:${HOMEPAGE:S,https://,,}",
		"/directory with spaces")

	vt.Output(
		"WARN: fname:1: All components of PATH (in this case \".\") should be absolute paths.",
		"WARN: fname:1: All components of PATH (in this case \"\") should be absolute paths.",
		"WARN: fname:2: \"/directory with spaces\" is not a valid pathname.")
}

func (s *Suite) Test_VartypeCheck_Pathmask(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Pathmask)

	vt.Varname("DISTDIRS")
	vt.Values(
		"/home/user/*",
		"src/*&*",
		"src/*/*")

	vt.Output(
		"WARN: fname:1: Found absolute pathname: /home/user/*",
		"WARN: fname:2: \"src/*&*\" is not a valid pathname mask.")

	vt.Op(opUseMatch)
	vt.Values("any")

	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_Pathname(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Pathname)

	vt.Varname("EGDIR")
	vt.Values(
		"${PREFIX}/*",
		"${PREFIX}/share/locale",
		"share/locale",
		"/bin")
	vt.Op(opUseMatch)
	vt.Values(
		"anything")

	vt.Output(
		"WARN: fname:1: \"${PREFIX}/*\" is not a valid pathname.",
		"WARN: fname:4: Found absolute pathname: /bin")
}

func (s *Suite) Test_VartypeCheck_Perl5Packlist(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Perl5Packlist)

	vt.Varname("PERL5_PACKLIST")
	vt.Values(
		"${PKGBASE}",
		"anything else")

	vt.Output(
		"WARN: fname:1: PERL5_PACKLIST should not depend on other variables.")
}

func (s *Suite) Test_VartypeCheck_Perms(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Perms)

	vt.Varname("CONF_FILES_PERMS")
	vt.Op(opAssignAppend)
	vt.Values(
		"root",
		"${ROOT_USER}",
		"ROOT_USER",
		"${REAL_ROOT_USER}")

	vt.Output(
		"ERROR: fname:2: ROOT_USER must not be used in permission definitions. Use REAL_ROOT_USER instead.")
}

func (s *Suite) Test_VartypeCheck_Pkgname(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).PkgName)

	vt.Varname("PKGNAME")
	vt.Values(
		"pkgbase-0",
		"pkgbase-1.0",
		"pkgbase-1.1234567890",
		"pkgbase-1z",
		"pkgbase-client-11a",
		"pkgbase-client-1.a",
		"pkgbase-client-1_20180101",
		"pkgbase-z1",
		"pkgbase-3.1.4.1.5.9.2.6.5.3.5.8.9.7.9")

	vt.Output(
		"WARN: fname:8: \"pkgbase-z1\" is not a valid package name. " +
			"A valid package name has the form packagename-version, " +
			"where version consists only of digits, letters and dots.")
}

func (s *Suite) Test_VartypeCheck_PkgOptionsVar(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).PkgOptionsVar)

	vt.Varname("PKG_OPTIONS_VAR.screen")
	vt.Values(
		"PKG_OPTIONS.${PKGBASE}",
		"PKG_OPTIONS.anypkgbase",
		"PKG_OPTS.mc")

	vt.Output(
		"ERROR: fname:1: PKGBASE must not be used in PKG_OPTIONS_VAR.",
		"ERROR: fname:3: PKG_OPTIONS_VAR must be of the form \"PKG_OPTIONS.*\", not \"PKG_OPTS.mc\".")
}

func (s *Suite) Test_VartypeCheck_PkgPath(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).PkgPath)

	t.CreateFileLines("category/other-package/Makefile")
	t.Chdir("category/package")

	vt.Varname("PKGPATH")
	vt.Values(
		"category/other-package",
		"${OTHER_VAR}",
		"invalid",
		"../../invalid/relative")

	vt.Output(
		"ERROR: fname:3: \"../../invalid\" does not exist.",
		"WARN: fname:3: \"../../invalid\" is not a valid relative package directory.",
		"ERROR: fname:4: \"../../../../invalid/relative\" does not exist.",
		"WARN: fname:4: \"../../../../invalid/relative\" is not a valid relative package directory.")
}

func (s *Suite) Test_VartypeCheck_PkgRevision(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).PkgRevision)

	vt.Varname("PKGREVISION")
	vt.Values(
		"3a")

	vt.Output(
		"WARN: fname:1: PKGREVISION must be a positive integer number.",
		"ERROR: fname:1: PKGREVISION only makes sense directly in the package Makefile.")

	vt.File("Makefile")
	vt.Values(
		"3")

	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_MachinePlatformPattern(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).MachinePlatformPattern)

	vt.Varname("ONLY_FOR_PLATFORM")
	vt.Op(opUseMatch)
	vt.Values(
		"linux-i386",
		"nextbsd-5.0-8087",
		"netbsd-7.0-l*",
		"NetBSD-1.6.2-i386",
		"FreeBSD*",
		"FreeBSD-*",
		"${LINUX}")

	vt.Output(
		"WARN: fname:1: \"linux-i386\" is not a valid platform pattern.",
		"WARN: fname:2: The pattern \"nextbsd\" cannot match any of "+
			"{ AIX BSDOS Bitrig Cygwin Darwin DragonFly FreeBSD FreeMiNT GNUkFreeBSD HPUX Haiku "+
			"IRIX Interix Linux Minix MirBSD NetBSD OSF1 OpenBSD QNX SCO_SV SunOS UnixWare "+
			"} for the operating system part of ONLY_FOR_PLATFORM.",
		"WARN: fname:2: The pattern \"8087\" cannot match any of "+
			"{ aarch64 aarch64eb alpha amd64 arc arm arm26 arm32 cobalt coldfire convex dreamcast "+
			"earm earmeb earmhf earmhfeb earmv4 earmv4eb earmv5 earmv5eb earmv6 earmv6eb earmv6hf "+
			"earmv6hfeb earmv7 earmv7eb earmv7hf earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 "+
			"i386 i586 i686 ia64 m68000 m68k m88k mips mips64 mips64eb mips64el mipseb mipsel mipsn32 "+
			"mlrisc ns32k pc532 pmax powerpc powerpc64 rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 "+
			"} for the hardware architecture part of ONLY_FOR_PLATFORM.",
		"WARN: fname:3: The pattern \"netbsd\" cannot match any of "+
			"{ AIX BSDOS Bitrig Cygwin Darwin DragonFly FreeBSD FreeMiNT GNUkFreeBSD HPUX Haiku "+
			"IRIX Interix Linux Minix MirBSD NetBSD OSF1 OpenBSD QNX SCO_SV SunOS UnixWare "+
			"} for the operating system part of ONLY_FOR_PLATFORM.",
		"WARN: fname:3: The pattern \"l*\" cannot match any of "+
			"{ aarch64 aarch64eb alpha amd64 arc arm arm26 arm32 cobalt coldfire convex dreamcast "+
			"earm earmeb earmhf earmhfeb earmv4 earmv4eb earmv5 earmv5eb earmv6 earmv6eb earmv6hf "+
			"earmv6hfeb earmv7 earmv7eb earmv7hf earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 "+
			"i386 i586 i686 ia64 m68000 m68k m88k mips mips64 mips64eb mips64el mipseb mipsel mipsn32 "+
			"mlrisc ns32k pc532 pmax powerpc powerpc64 rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 "+
			"} for the hardware architecture part of ONLY_FOR_PLATFORM.",
		"WARN: fname:5: \"FreeBSD*\" is not a valid platform pattern.")
}

func (s *Suite) Test_VartypeCheck_PythonDependency(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).PythonDependency)

	vt.Varname("PYTHON_VERSIONED_DEPENDENCIES")
	vt.Values(
		"cairo",
		"${PYDEP}",
		"cairo,X")

	vt.Output(
		"WARN: fname:2: Python dependencies should not contain variables.",
		"WARN: fname:3: Invalid Python dependency \"cairo,X\".")
}

func (s *Suite) Test_VartypeCheck_PrefixPathname(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).PrefixPathname)

	vt.Varname("PKGMANDIR")
	vt.Values(
		"man/man1",
		"share/locale")

	vt.Output(
		"WARN: fname:1: Please use \"${PKGMANDIR}/man1\" instead of \"man/man1\".")
}

func (s *Suite) Test_VartypeCheck_RelativePkgPath(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).RelativePkgPath)

	t.CreateFileLines("category/other-package/Makefile")
	t.Chdir("category/package")

	vt.Varname("DISTINFO_FILE")
	vt.Values(
		"category/other-package",
		"../../category/other-package",
		"${OTHER_VAR}",
		"invalid",
		"../../invalid/relative")

	vt.Output(
		"ERROR: fname:1: \"category/other-package\" does not exist.",
		"ERROR: fname:4: \"invalid\" does not exist.",
		"ERROR: fname:5: \"../../invalid/relative\" does not exist.")
}

func (s *Suite) Test_VartypeCheck_Restricted(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Restricted)

	vt.Varname("NO_BIN_ON_CDROM")
	vt.Values(
		"May only be distributed free of charge")

	vt.Output(
		"WARN: fname:1: The only valid value for NO_BIN_ON_CDROM is ${RESTRICTED}.")
}

func (s *Suite) Test_VartypeCheck_SedCommands(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).SedCommands)

	vt.Varname("SUBST_SED.dummy")
	vt.Values(
		"s,@COMPILER@,gcc,g",
		"-e s,a,b, -e a,b,c,",
		"-e \"s,#,comment ,\"",
		"-e \"s,\\#,comment ,\"",
		"-E",
		"-n",
		"-e 1d",
		"1d",
		"-e")

	vt.Output(
		"NOTE: fname:1: Please always use \"-e\" in sed commands, even if there is only one substitution.",
		"NOTE: fname:2: Each sed command should appear in an assignment of its own.",
		"WARN: fname:3: The # character starts a comment.",
		"ERROR: fname:3: Invalid shell words \"\\\"s,\" in sed commands.",
		"WARN: fname:8: Unknown sed command \"1d\".",
		"ERROR: fname:9: The -e option to sed requires an argument.")
}

func (s *Suite) Test_VartypeCheck_ShellCommand(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).ShellCommand)

	vt.Varname("INSTALL_CMD")
	vt.Values(
		"${INSTALL_DATA} -m 0644 ${WRKDIR}/source ${DESTDIR}${PREFIX}/target")

	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_ShellCommands(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).ShellCommands)

	vt.Varname("GENERATE_PLIST")
	vt.Values(
		"echo bin/program",
		"echo bin/program;")

	vt.Output(
		"WARN: fname:1: This shell command list should end with a semicolon.")
}

func (s *Suite) Test_VartypeCheck_Stage(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Stage)

	vt.Varname("SUBST_STAGE.dummy")
	vt.Values(
		"post-patch",
		"post-modern",
		"pre-test")

	vt.Output(
		"WARN: fname:2: Invalid stage name \"post-modern\". " +
			"Use one of {pre,do,post}-{extract,patch,configure,build,test,install}.")
}

func (s *Suite) Test_VartypeCheck_Tool(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).Tool)

	t.SetupToolUsable("tool1", "")
	t.SetupToolUsable("tool2", "")
	t.SetupToolUsable("tool3", "")

	vt.Varname("USE_TOOLS")
	vt.Op(opAssignAppend)
	vt.Values(
		"tool3:run",
		"tool2:unknown",
		"${t}",
		"mal:formed:tool",
		"unknown")

	vt.Output(
		"ERROR: fname:2: Unknown tool dependency \"unknown\". "+
			"Use one of \"bootstrap\", \"build\", \"pkgsrc\", \"run\" or \"test\".",
		"ERROR: fname:4: Malformed tool dependency: \"mal:formed:tool\".",
		"ERROR: fname:5: Unknown tool \"unknown\".")

	vt.Varname("USE_TOOLS.NetBSD")
	vt.Op(opAssignAppend)
	vt.Values(
		"tool3:run",
		"tool2:unknown")

	vt.Output(
		"ERROR: fname:12: Unknown tool dependency \"unknown\". " +
			"Use one of \"bootstrap\", \"build\", \"pkgsrc\", \"run\" or \"test\".")

	vt.Varname("TOOLS_NOOP")
	vt.Op(opAssignAppend)
	vt.Values(
		"gmake:run")

	vt.OutputEmpty()
}

func (s *Suite) Test_VartypeCheck_URL(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).URL)

	vt.Varname("HOMEPAGE")
	vt.Values(
		"# none",
		"${OTHER_VAR}",
		"https://www.netbsd.org/",
		"mailto:someone@example.org",
		"httpxs://www.example.org",
		"https://www.example.org",
		"https://www.example.org/path with spaces",
		"string with spaces")

	vt.Output(
		"WARN: fname:3: Please write NetBSD.org instead of www.netbsd.org.",
		"WARN: fname:4: \"mailto:someone@example.org\" is not a valid URL.",
		"WARN: fname:5: \"httpxs://www.example.org\" is not a valid URL. Only ftp, gopher, http, and https URLs are allowed here.",
		"NOTE: fname:6: For consistency, please add a trailing slash to \"https://www.example.org\".",
		"WARN: fname:7: \"https://www.example.org/path with spaces\" is not a valid URL.",
		"WARN: fname:8: \"string with spaces\" is not a valid URL.")
}

func (s *Suite) Test_VartypeCheck_UserGroupName(c *check.C) {
	t := s.Init(c)
	vt := NewVartypeCheckTester(t, (*VartypeCheck).UserGroupName)

	vt.Varname("APACHE_USER")
	vt.Values(
		"user with spaces",
		"typical_username",
		"user123",
		"domain\\user",
		"${OTHER_VAR}")

	vt.Output(
		"WARN: fname:1: Invalid user or group name \"user with spaces\".",
		"WARN: fname:4: Invalid user or group name \"domain\\\\user\".")
}

func (s *Suite) Test_VartypeCheck_VariableName(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).VariableName)

	vt.Varname("BUILD_DEFS")
	vt.Values(
		"VARBASE",
		"VarBase",
		"PKG_OPTIONS_VAR.pkgbase",
		"${INDIRECT}")

	vt.Output(
		"WARN: fname:2: \"VarBase\" is not a valid variable name.")
}

func (s *Suite) Test_VartypeCheck_Version(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Version)

	vt.Varname("PERL5_REQD")
	vt.Op(opAssignAppend)
	vt.Values(
		"0",
		"1.2.3.4.5.6",
		"4.1nb17",
		"4.1-SNAPSHOT",
		"4pre7")
	vt.Output(
		"WARN: fname:4: Invalid version number \"4.1-SNAPSHOT\".")

	vt.Op(opUseMatch)
	vt.Values(
		"a*",
		"1.2/456",
		"[0-9]*")
	vt.Output(
		"WARN: fname:11: Invalid version number pattern \"a*\".",
		"WARN: fname:12: Invalid version number pattern \"1.2/456\".")
}

func (s *Suite) Test_VartypeCheck_WrapperReorder(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).WrapperReorder)

	vt.Varname("WRAPPER_REORDER")
	vt.Op(opAssignAppend)
	vt.Values(
		"reorder:l:first:second",
		"reorder:l:first",
		"omit:first")
	vt.Output(
		"WARN: fname:2: Unknown wrapper reorder command \"reorder:l:first\".",
		"WARN: fname:3: Unknown wrapper reorder command \"omit:first\".")
}

func (s *Suite) Test_VartypeCheck_WrapperTransform(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).WrapperTransform)

	vt.Varname("WRAPPER_TRANSFORM")
	vt.Op(opAssignAppend)
	vt.Values(
		"rm:-O3",
		"opt:-option",
		"rename:src:dst",
		"rm-optarg:-option",
		"rmdir:/usr/include",
		"rpath:/usr/lib:/usr/pkg/lib",
		"rpath:/usr/lib",
		"unknown")
	vt.Output(
		"WARN: fname:7: Unknown wrapper transform command \"rpath:/usr/lib\".",
		"WARN: fname:8: Unknown wrapper transform command \"unknown\".")
}

func (s *Suite) Test_VartypeCheck_WrksrcSubdirectory(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).WrksrcSubdirectory)

	vt.Varname("BUILD_DIRS")
	vt.Op(opAssignAppend)
	vt.Values(
		"${WRKSRC}",
		"${WRKSRC}/",
		"${WRKSRC}/.",
		"${WRKSRC}/subdir",
		"${CONFIGURE_DIRS}",
		"${WRKSRC}/directory with spaces",
		"directory with spaces")
	vt.Output(
		"NOTE: fname:1: You can use \".\" instead of \"${WRKSRC}\".",
		"NOTE: fname:2: You can use \".\" instead of \"${WRKSRC}/\".",
		"NOTE: fname:3: You can use \".\" instead of \"${WRKSRC}/.\".",
		"NOTE: fname:4: You can use \"subdir\" instead of \"${WRKSRC}/subdir\".",
		"NOTE: fname:6: You can use \"directory with spaces\" instead of \"${WRKSRC}/directory with spaces\".",
		"WARN: fname:7: \"directory with spaces\" is not a valid subdirectory of ${WRKSRC}.")
}

func (s *Suite) Test_VartypeCheck_Yes(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).Yes)

	vt.Varname("APACHE_MODULE")
	vt.Values(
		"yes",
		"no",
		"${YESVAR}")

	vt.Output(
		"WARN: fname:2: APACHE_MODULE should be set to YES or yes.",
		"WARN: fname:3: APACHE_MODULE should be set to YES or yes.")

	vt.Varname("PKG_DEVELOPER")
	vt.Op(opUseMatch)
	vt.Values(
		"yes",
		"no",
		"${YESVAR}")

	vt.Output(
		"WARN: fname:11: PKG_DEVELOPER should only be used in a \".if defined(...)\" condition.",
		"WARN: fname:12: PKG_DEVELOPER should only be used in a \".if defined(...)\" condition.",
		"WARN: fname:13: PKG_DEVELOPER should only be used in a \".if defined(...)\" condition.")
}

func (s *Suite) Test_VartypeCheck_YesNo(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).YesNo)

	vt.Varname("GNU_CONFIGURE")
	vt.Values(
		"yes",
		"no",
		"ja",
		"${YESVAR}")

	vt.Output(
		"WARN: fname:3: GNU_CONFIGURE should be set to YES, yes, NO, or no.",
		"WARN: fname:4: GNU_CONFIGURE should be set to YES, yes, NO, or no.")
}

func (s *Suite) Test_VartypeCheck_YesNoIndirectly(c *check.C) {
	vt := NewVartypeCheckTester(s.Init(c), (*VartypeCheck).YesNoIndirectly)

	vt.Varname("GNU_CONFIGURE")
	vt.Values(
		"yes",
		"no",
		"ja",
		"${YESVAR}")

	vt.Output(
		"WARN: fname:3: GNU_CONFIGURE should be set to YES, yes, NO, or no.")
}

// VartypeCheckTester helps to test the many different checks in VartypeCheck.
// It keeps track of the current variable, operator, file name, line number,
// so that the test can focus on the interesting details.
type VartypeCheckTester struct {
	tester   *Tester
	checker  func(cv *VartypeCheck)
	fileName string
	lineno   int
	varname  string
	op       MkOperator
}

// NewVartypeCheckTester starts the test with a file name of "fname", at line 1,
// with "=" as the operator. The variable has to be initialized explicitly.
func NewVartypeCheckTester(t *Tester, checker func(cv *VartypeCheck)) *VartypeCheckTester {
	return &VartypeCheckTester{
		t,
		checker,
		"fname",
		1,
		"",
		opAssign}
}

func (vt *VartypeCheckTester) Varname(varname string) {
	vt.varname = varname
	vt.nextSection()
}

func (vt *VartypeCheckTester) File(fileName string) {
	vt.fileName = fileName
	vt.lineno = 1
}

// Op sets the operator for the following tests.
// The line number is advanced to the next number ending in 1, e.g. 11, 21, 31.
func (vt *VartypeCheckTester) Op(op MkOperator) {
	vt.op = op
	vt.nextSection()
}

// Values feeds each of the values to the actual check.
// Each value is interpreted as if it were written verbatim into a Makefile line.
// That is, # starts a comment, and for the opUseMatch operator, all closing braces must be escaped.
func (vt *VartypeCheckTester) Values(values ...string) {
	for _, value := range values {
		op := vt.op
		opStr := op.String()
		varname := vt.varname

		var text string
		switch {
		case contains(opStr, "="):
			if hasSuffix(varname, "+") && opStr == "=" {
				text = varname + " " + opStr + value
			} else {
				text = varname + opStr + value
			}
		case op == opUseMatch:
			text = fmt.Sprintf(".if ${%s:M%s} == \"\"", varname, value)
		default:
			panic("Invalid operator: " + opStr)
		}

		mkline := vt.tester.NewMkLine(vt.fileName, vt.lineno, text)
		comment := ""
		if mkline.IsVarassign() {
			mkline.Tokenize(value, true) // Produce some warnings as side-effects.
			comment = mkline.VarassignComment()
		}

		effectiveValue := value
		if mkline.IsVarassign() {
			effectiveValue = mkline.Value()
		}

		valueNovar := mkline.WithoutMakeVariables(effectiveValue)
		vc := &VartypeCheck{mkline, mkline.Line, varname, op, effectiveValue, valueNovar, comment, false}
		vt.checker(vc)

		vt.lineno++
	}
}

// Output checks that the output from all previous steps is
// the same as the expectedLines.
func (vt *VartypeCheckTester) Output(expectedLines ...string) {
	vt.tester.CheckOutputLines(expectedLines...)
}

func (vt *VartypeCheckTester) OutputEmpty() {
	vt.tester.CheckOutputEmpty()
}

func (vt *VartypeCheckTester) nextSection() {
	if vt.lineno%10 != 1 {
		vt.lineno = vt.lineno - vt.lineno%10 + 11
	}
}
