package main

import "gopkg.in/check.v1"

func (s *Suite) Test_MkLineChecker_CheckVartype__simple_type(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wtypes")
	t.SetupVartypes()

	vartype1 := G.Pkgsrc.vartypes["COMMENT"]
	c.Assert(vartype1, check.NotNil)
	c.Check(vartype1.guessed, equals, false)

	vartype := G.Pkgsrc.VariableType("COMMENT")

	c.Assert(vartype, check.NotNil)
	c.Check(vartype.basicType.name, equals, "Comment")
	c.Check(vartype.guessed, equals, false)
	c.Check(vartype.kindOfList, equals, lkNone)

	MkLineChecker{dummyMkLine}.CheckVartype("COMMENT", opAssign, "A nice package", "")

	t.CheckOutputLines(
		"WARN: COMMENT should not begin with \"A\".")
}

func (s *Suite) Test_MkLineChecker_CheckVartype(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mkline := t.NewMkLine("fname", 1, "DISTNAME=gcc-${GCC_VERSION}")

	MkLineChecker{mkline}.CheckVartype("DISTNAME", opAssign, "gcc-${GCC_VERSION}", "")

	t.CheckOutputEmpty()
}

// Pkglint once interpreted all lists as consisting of shell tokens,
// splitting this URL at the ampersand.
func (s *Suite) Test_MkLineChecker_checkVarassign__URL_with_shell_special_characters(c *check.C) {
	t := s.Init(c)

	G.Pkg = NewPackage(t.File("graphics/gimp-fix-ca"))
	t.SetupVartypes()
	mkline := t.NewMkLine("fname", 10, "MASTER_SITES=http://registry.gimp.org/file/fix-ca.c?action=download&id=9884&file=")

	MkLineChecker{mkline}.checkVarassign()

	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_Check__conditions(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wtypes")
	t.SetupVartypes()

	testCond := func(cond string, output ...string) {
		MkLineChecker{t.NewMkLine("fname", 1, cond)}.checkDirectiveCond()
		t.CheckOutputLines(output...)
	}

	testCond(".if !empty(PKGSRC_COMPILER:Mmycc)",
		"WARN: fname:1: The pattern \"mycc\" cannot match any of "+
			"{ ccache ccc clang distcc f2c gcc hp icc ido "+
			"mipspro mipspro-ucode pcc sunpro xlc } for PKGSRC_COMPILER.")

	testCond(".elif ${A} != ${B}")

	testCond(".if ${HOMEPAGE} == \"mailto:someone@example.org\"",
		"WARN: fname:1: \"mailto:someone@example.org\" is not a valid URL.")

	testCond(".if !empty(PKGSRC_RUN_TEST:M[Y][eE][sS])",
		"WARN: fname:1: PKGSRC_RUN_TEST should be matched against \"[yY][eE][sS]\" or \"[nN][oO]\", not \"[Y][eE][sS]\".")

	testCond(".if !empty(IS_BUILTIN.Xfixes:M[yY][eE][sS])")

	testCond(".if !empty(${IS_BUILTIN.Xfixes:M[yY][eE][sS]})",
		"WARN: fname:1: The empty() function takes a variable name as parameter, not a variable expression.")

	testCond(".if ${EMUL_PLATFORM} == \"linux-x386\"",
		"WARN: fname:1: "+
			"\"x386\" is not valid for the hardware architecture part of EMUL_PLATFORM. "+
			"Use one of "+
			"{ aarch64 aarch64eb alpha amd64 arc arm arm26 arm32 cobalt coldfire convex "+
			"dreamcast earm earmeb earmhf earmhfeb earmv4 earmv4eb earmv5 earmv5eb earmv6 earmv6eb "+
			"earmv6hf earmv6hfeb earmv7 earmv7eb earmv7hf earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 "+
			"i386 i586 i686 ia64 m68000 m68k m88k mips mips64 mips64eb mips64el mipseb mipsel mipsn32 "+
			"mlrisc ns32k pc532 pmax powerpc powerpc64 rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 "+
			"} instead.")

	testCond(".if ${EMUL_PLATFORM:Mlinux-x386}",
		"WARN: fname:1: "+
			"The pattern \"x386\" cannot match any of { aarch64 aarch64eb alpha amd64 arc arm arm26 "+
			"arm32 cobalt coldfire convex dreamcast earm earmeb earmhf earmhfeb earmv4 earmv4eb "+
			"earmv5 earmv5eb earmv6 earmv6eb earmv6hf earmv6hfeb earmv7 earmv7eb earmv7hf "+
			"earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 i386 i586 i686 ia64 m68000 m68k m88k "+
			"mips mips64 mips64eb mips64el mipseb mipsel mipsn32 mlrisc ns32k pc532 pmax powerpc powerpc64 "+
			"rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 } "+
			"for the hardware architecture part of EMUL_PLATFORM.",
		"NOTE: fname:1: EMUL_PLATFORM should be compared using == instead of the :M or :N modifier without wildcards.")

	testCond(".if ${MACHINE_PLATFORM:MUnknownOS-*-*} || ${MACHINE_ARCH:Mx86}",
		"WARN: fname:1: "+
			"The pattern \"UnknownOS\" cannot match any of "+
			"{ AIX BSDOS Bitrig Cygwin Darwin DragonFly FreeBSD FreeMiNT GNUkFreeBSD HPUX Haiku "+
			"IRIX Interix Linux Minix MirBSD NetBSD OSF1 OpenBSD QNX SCO_SV SunOS UnixWare "+
			"} for the operating system part of MACHINE_PLATFORM.",
		"WARN: fname:1: "+
			"The pattern \"x86\" cannot match any of "+
			"{ aarch64 aarch64eb alpha amd64 arc arm arm26 arm32 cobalt coldfire convex dreamcast earm "+
			"earmeb earmhf earmhfeb earmv4 earmv4eb earmv5 earmv5eb earmv6 earmv6eb earmv6hf earmv6hfeb "+
			"earmv7 earmv7eb earmv7hf earmv7hfeb evbarm hpcmips hpcsh hppa hppa64 i386 i586 i686 ia64 "+
			"m68000 m68k m88k mips mips64 mips64eb mips64el mipseb mipsel mipsn32 mlrisc ns32k pc532 pmax "+
			"powerpc powerpc64 rs6000 s390 sh3eb sh3el sparc sparc64 vax x86_64 "+
			"} for MACHINE_ARCH.",
		"NOTE: fname:1: MACHINE_ARCH should be compared using == instead of the :M or :N modifier without wildcards.")

	testCond(".if ${MASTER_SITES:Mftp://*} == \"ftp://netbsd.org/\"")
}

func (s *Suite) Test_MkLineChecker_checkVarassign(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()

	G.Mk = t.NewMkLines("Makefile",
		MkRcsID,
		"ac_cv_libpari_libs+=\t-L${BUILDLINK_PREFIX.pari}/lib") // From math/clisp-pari/Makefile, rev. 1.8

	MkLineChecker{G.Mk.mklines[1]}.checkVarassign()

	t.CheckOutputLines(
		"WARN: Makefile:2: ac_cv_libpari_libs is defined but not used.")
}

func (s *Suite) Test_MkLineChecker_checkVarassignPermissions(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	mkline := t.NewMkLine("options.mk", 2, "PKG_DEVELOPER?=\tyes")

	MkLineChecker{mkline}.checkVarassignPermissions()

	t.CheckOutputLines(
		"WARN: options.mk:2: The variable PKG_DEVELOPER may not be given a default value by any package.")
}

// Don't check the permissions for infrastructure files since they have their own rules.
func (s *Suite) Test_MkLineChecker_checkVarassignDefPermissions__infrastructure(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	t.SetupFileMkLines("mk/infra.mk",
		MkRcsID,
		"",
		"PKG_DEVELOPER?=\tyes")
	t.SetupFileMkLines("mk/bsd.pkg.mk")

	G.CheckDirent(t.File("mk/infra.mk"))

	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_CheckVarusePermissions(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	mklines := t.NewMkLines("options.mk",
		MkRcsID,
		"COMMENT=\t${GAMES_USER}",
		"COMMENT:=\t${PKGBASE}",
		"PYPKGPREFIX=${PKGBASE}")
	G.Pkgsrc.UserDefinedVars.Define("GAMES_USER", mklines.mklines[0])

	mklines.Check()

	t.CheckOutputLines(
		"WARN: options.mk:3: PKGBASE should not be evaluated at load time.",
		"WARN: options.mk:4: The variable PYPKGPREFIX may not be set in this file; it would be ok in pyversion.mk.",
		"WARN: options.mk:4: PKGBASE should not be evaluated indirectly at load time.",
		"NOTE: options.mk:4: This variable value should be aligned to column 17.")
}

func (s *Suite) Test_MkLineChecker_CheckVarusePermissions__load_time(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	mklines := t.NewMkLines("options.mk",
		MkRcsID,
		"WRKSRC:=${.CURDIR}")

	mklines.Check()

	// Don't warn that ".CURDIR should not be evaluated at load time."
	t.CheckOutputLines(
		"NOTE: options.mk:2: This variable value should be aligned to column 17.")
}

func (s *Suite) Test_MkLineChecker_WarnVaruseLocalbase(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("options.mk", 56, "PKGNAME=${LOCALBASE}")

	MkLineChecker{mkline}.WarnVaruseLocalbase()

	t.CheckOutputLines(
		"WARN: options.mk:56: Please use PREFIX instead of LOCALBASE.")
}

func (s *Suite) Test_MkLineChecker_CheckRelativePkgdir(c *check.C) {
	t := s.Init(c)

	mkline := t.NewMkLine("Makefile", 46, "# dummy")

	MkLineChecker{mkline}.CheckRelativePkgdir("../pkgbase")

	t.CheckOutputLines(
		"ERROR: Makefile:46: \"../pkgbase\" does not exist.",
		"WARN: Makefile:46: \"../pkgbase\" is not a valid relative package directory.")
}

// PR pkg/46570, item 2
func (s *Suite) Test_MkLineChecker__unclosed_varuse(c *check.C) {
	t := s.Init(c)

	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"EGDIRS=${EGDIR/apparmor.d ${EGDIR/dbus-1/system.d ${EGDIR/pam.d")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: Unclosed Make variable starting at \"${EGDIR/apparmor.d $...\".",
		"WARN: Makefile:2: EGDIRS is defined but not used.",
		"WARN: Makefile:2: Pkglint parse error in MkLine.Tokenize at "+
			"\"${EGDIR/apparmor.d ${EGDIR/dbus-1/system.d ${EGDIR/pam.d\".")
}

func (s *Suite) Test_MkLineChecker__Varuse_Modifier_L(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	G.Mk = t.NewMkLines("x11/xkeyboard-config/Makefile",
		"FILES_SUBST+=XKBCOMP_SYMLINK=${${XKBBASE}/xkbcomp:L:Q}")

	MkLineChecker{G.Mk.mklines[0]}.Check()

	// Don't warn that ${XKBBASE}/xkbcomp is used but not defined.
	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_CheckCond__comparison_with_shell_command(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	G.Mk = t.NewMkLines("security/openssl/Makefile",
		MkRcsID,
		".if ${PKGSRC_COMPILER} == \"gcc\" && ${CC} == \"cc\"",
		".endif")

	G.Mk.Check()

	// Don't warn about unknown shell command "cc".
	t.CheckOutputLines(
		"WARN: security/openssl/Makefile:2: Use ${PKGSRC_COMPILER:Mgcc} instead of the == operator.")
}

func (s *Suite) Test_MkLine_CheckCond_comparing_PKGSRC_COMPILER_with_eqeq(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	G.Mk = t.NewMkLines("audio/pulseaudio/Makefile",
		MkRcsID,
		".if ${OPSYS} == \"Darwin\" && ${PKGSRC_COMPILER} == \"clang\"",
		".endif")

	G.Mk.Check()

	t.CheckOutputLines(
		"WARN: audio/pulseaudio/Makefile:2: Use ${PKGSRC_COMPILER:Mclang} instead of the == operator.")
}

func (s *Suite) Test_MkLineChecker_CheckVartype__CFLAGS_with_backticks(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	G.Mk = t.NewMkLines("chat/pidgin-icb/Makefile",
		MkRcsID,
		"CFLAGS+=\t`pkg-config pidgin --cflags`")
	mkline := G.Mk.mklines[1]

	words, rest := splitIntoMkWords(mkline.Line, mkline.Value())

	c.Check(words, deepEquals, []string{"`pkg-config pidgin --cflags`"})
	c.Check(rest, equals, "")

	MkLineChecker{G.Mk.mklines[1]}.CheckVartype("CFLAGS", opAssignAppend, "`pkg-config pidgin --cflags`", "")

	// No warning about "`pkg-config" being an unknown CFlag.
	t.CheckOutputEmpty()
}

// See PR 46570, Ctrl+F "4. Shell quoting".
// Pkglint is correct, since the shell sees this definition for
// CPPFLAGS as three words, not one word.
func (s *Suite) Test_MkLineChecker_CheckVartype_CFLAGS(c *check.C) {
	t := s.Init(c)

	t.SetupVartypes()
	mklines := t.NewMkLines("Makefile",
		MkRcsID,
		"CPPFLAGS.SunOS+=\t-DPIPECOMMAND=\\\"/usr/sbin/sendmail -bs %s\\\"")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: Makefile:2: Unknown compiler flag \"-bs\".",
		"WARN: Makefile:2: Compiler flag \"%s\\\\\\\"\" should start with a hyphen.")
}

// Up to 2018-01-28, pkglint applied the autofix also to the continuation
// lines, which is incorrect. It replaced the dot in "4.*" with spaces.
func (s *Suite) Test_MkLineChecker_checkDirectiveIndentation_autofix(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall", "--autofix")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		".if ${PKGNAME} == pkgname",
		".if \\",
		"   ${PLATFORM:MNetBSD-4.*}",
		".endif",
		".endif")

	mklines.Check()

	t.CheckOutputLines(
		"AUTOFIX: ~/options.mk:3: Replacing \".\" with \".  \".",
		"AUTOFIX: ~/options.mk:5: Replacing \".\" with \".  \".")

	t.CheckFileLines("options.mk",
		MkRcsID,
		".if ${PKGNAME} == pkgname",
		".  if \\",
		"   ${PLATFORM:MNetBSD-4.*}",
		".  endif",
		".endif")
}

func (s *Suite) Test_MkLineChecker_CheckVaruseShellword(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"GOPATH=\t${WRKDIR}",
		"do-build:",
		"\tcd ${WRKSRC} && GOPATH=${GOPATH} PATH=${PATH} :")

	mklines.Check()

	// For WRKSRC and GOPATH, no quoting is necessary since pkgsrc directories by
	// definition don't contain special characters. Therefore they don't need the
	// :Q, not even when used as part of a shell word.

	// For PATH, the quoting is necessary because it may contain directories outside
	// of pkgsrc, and these may contain special characters.

	t.CheckOutputLines(
		"WARN: ~/options.mk:4: The variable PATH should be quoted as part of a shell word.")
}

// The ${VARNAME:=suffix} expression should only be used with lists.
// It typically appears in MASTER_SITE definitions.
func (s *Suite) Test_MkLineChecker_CheckVaruse__eq_nonlist(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	t.SetupMasterSite("MASTER_SITE_GITHUB", "https://github.com/")
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"WRKSRC=\t\t${WRKDIR:=/subdir}",
		"MASTER_SITES=\t${MASTER_SITE_GITHUB:=organization/}")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/options.mk:2: The :from=to modifier should only be used with lists.")
}

func (s *Suite) Test_MkLineChecker_CheckVaruse__for(c *check.C) {
	t := s.Init(c)

	t.SetupCommandLine("-Wall")
	t.SetupVartypes()
	t.SetupMasterSite("MASTER_SITE_GITHUB", "https://github.com/")
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		".for var in a b c",
		"\t: ${var}",
		".endfor")

	mklines.Check()

	t.CheckOutputEmpty()
}

func (s *Suite) Test_MkLineChecker_CheckVaruse__build_defs(c *check.C) {
	t := s.Init(c)

	// XXX: This paragraph should not be necessary since VARBASE and X11_TYPE
	// are also defined in vardefs.go.
	t.SetupPkgsrc()
	t.CreateFileLines("mk/defaults/mk.conf",
		"VARBASE?= /usr/pkg/var")
	G.Pkgsrc.LoadInfrastructure()

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("options.mk",
		MkRcsID,
		"COMMENT=        ${VARBASE} ${X11_TYPE}",
		"BUILD_DEFS+=    X11_TYPE")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/options.mk:2: The user-defined variable VARBASE is used but not added to BUILD_DEFS.")
}

func (s *Suite) Test_MkLineChecker_checkVarassignSpecific(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	G.Pkgsrc.LoadInfrastructure()

	t.SetupCommandLine("-Wall,no-space")
	t.SetupVartypes()
	mklines := t.SetupFileMkLines("module.mk",
		MkRcsID,
		"EGDIR=                  ${PREFIX}/etc/rc.d",
		"_TOOLS_VARNAME.sed=     SED",
		"DIST_SUBDIR=            ${PKGNAME}",
		"WRKSRC=                 ${PKGNAME}",
		"SITES_distfile.tar.gz=  ${MASTER_SITES_GITHUB:=user/}")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/module.mk:2: Please use the RCD_SCRIPTS mechanism to install rc.d scripts automatically to ${RCD_SCRIPTS_EXAMPLEDIR}.",
		"WARN: ~/module.mk:3: _TOOLS_VARNAME.sed is defined but not used.",
		"WARN: ~/module.mk:3: Variable names starting with an underscore (_TOOLS_VARNAME.sed) are reserved for internal pkgsrc use.",
		"WARN: ~/module.mk:4: PKGNAME should not be used in DIST_SUBDIR, as it includes the PKGREVISION. Please use PKGNAME_NOREV instead.",
		"WARN: ~/module.mk:5: PKGNAME should not be used in WRKSRC, as it includes the PKGREVISION. Please use PKGNAME_NOREV instead.",
		"WARN: ~/module.mk:6: SITES_distfile.tar.gz is defined but not used.",
		"WARN: ~/module.mk:6: SITES_* is deprecated. Please use SITES.* instead.")
}

func (s *Suite) Test_MkLineChecker_checkText(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	G.Pkgsrc.LoadInfrastructure()

	t.SetupCommandLine("-Wall,no-space")
	mklines := t.SetupFileMkLines("module.mk",
		MkRcsID,
		"CFLAGS+=                -Wl,--rpath,${PREFIX}/lib",
		"PKG_FAIL_REASON+=       \"Group ${GAMEGRP} doesn't exist.\"")

	mklines.Check()

	t.CheckOutputLines(
		"WARN: ~/module.mk:2: Please use ${COMPILER_RPATH_FLAG} instead of \"-Wl,--rpath,\".",
		"WARN: ~/module.mk:3: Use of \"GAMEGRP\" is deprecated. Use GAMES_GROUP instead.")
}

func (s *Suite) Test_MkLineChecker_CheckRelativePath(c *check.C) {
	t := s.Init(c)

	t.SetupPkgsrc()
	G.Pkgsrc.LoadInfrastructure()
	t.CreateFileLines("wip/package/Makefile")
	t.CreateFileLines("wip/package/module.mk")
	mklines := t.SetupFileMkLines("category/package/module.mk",
		MkRcsID,
		"DEPENDS+=       wip-package-[0-9]*:../../wip/package",
		".include \"../../wip/package/module.mk\"",
		"",
		"DEPENDS+=       unresolvable-[0-9]*:../../lang/${LATEST_PYTHON}",
		".include \"../../lang/${LATEST_PYTHON}/module.mk\"",
		"",
		".include \"module.mk\"",
		".include \"../../category/../category/package/module.mk\"", // Oops
		".include \"../../mk/bsd.prefs.mk\"",
		".include \"../package/module.mk\"")

	mklines.Check()

	t.CheckOutputLines(
		"ERROR: ~/category/package/module.mk:2: A main pkgsrc package must not depend on a pkgsrc-wip package.",
		"ERROR: ~/category/package/module.mk:3: A main pkgsrc package must not depend on a pkgsrc-wip package.",
		"WARN: ~/category/package/module.mk:11: Invalid relative path \"../package/module.mk\".")
}
