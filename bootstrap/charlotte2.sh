#!/bin/sh
#
# This is Charlotte's super-simplifed version of pkgsrc's 'bootstrap'
# script specifically for FreeMiNT.
#

set -e
set -x

# Don't let the bootstrap program get confused by a pre-existing mk.conf
# file.
MAKECONF=/dev/null
export MAKECONF

###############

user=charlotte
group=users
root_user=root
root_group=root
prefix=/e/pkg
make_jobs=1

unprivileged=yes
opsys=FreeMiNT
machine_arch=m68k

###############

# where the building takes place
bootstrapdir=`dirname "$0"`
bootstrapdir=`cd "${bootstrapdir}" && pwd`
pkgsrcdir=`dirname "${bootstrapdir}"`
wrkdir="`pwd`/work"

[ -z "$varbase" ] && varbase=${prefix}/var
[ -z "$pkgdbdir" ] && pkgdbdir=${prefix}/pkgdb
[ -z "$sysconfdir" ] && sysconfdir=${prefix}/etc
[ -z "$pkginfodir" ] && pkginfodir=info
[ -z "$pkgmandir" ] && pkgmandir=man
[ -z "$sysconfbase" ] && sysconfbase=${prefix}/etc
infodir=${prefix}/${pkginfodir}
mandir=${prefix}/${pkgmandir}
[ -z "$sysconfdir" ] && sysconfdir=${prefix}/etc

shprog=sh
sedprog=sed

need_bsd_install=no
configure_quiet_flags=""
make_quiet_flags=""

########################################################## 

die()
{
	echo >&2 "$*"
	exit 1
}

echo_msg()
{
	echo "===> $*"
}

# run a command, abort if it fails
run_cmd()
{
	echo_msg "running: $*"
	eval "$@"
	ret=$?
        if [ $ret -ne 0 ]; then
		echo_msg "exited with status $ret"
		die "aborted."
	fi
}

copy_src()
{
	_src="$1"; _dst="$2"
	if [ ! -d $wrkdir/$_dst ]; then
		mkdir -p $wrkdir/$_dst
	fi
	cp -r $_src/* $wrkdir/$_dst
	if [ -f $wrkdir/$_dst/config.guess ]; then
		cp $pkgsrcdir/mk/gnu-config/config.guess $wrkdir/$_dst/
	fi
	if [ -f $wrkdir/$_dst/config.sub ]; then
		cp $pkgsrcdir/mk/gnu-config/config.sub $wrkdir/$_dst/
	fi
}

########################################################## 

# MAIN:

# install-sh is going into bin
mkdir -p ${wrkdir}/bin

############################# build install-sh ####

build_install_sh() {
  run_cmd "sed -e 's|@DEFAULT_INSTALL_MODE@|'${default_install_mode-0755}'|' $pkgsrcdir/sysutils/install-sh/files/install-sh.in > $wrkdir/bin/install-sh"
  run_cmd "chmod +x $wrkdir/bin/install-sh"
  install_sh="sh $wrkdir/bin/install-sh"
  
  if [ $unprivileged = "yes" ]; then
  	echo_msg "building as unprivileged user $user/$group"
  
  	# force bmake install target to use $user and $group
  	echo "BINOWN=$user
  BINGRP=$group
  LIBOWN=$user
  LIBGRP=$group
  MANOWN=$user
  MANGRP=$group" > ${wrkdir}/Makefile.inc
  elif is_root; then
  	user=$root_user
  	group=$root_group
  else
  	die "You must be either root to install bootstrap-pkgsrc or use the --unprivileged option."
  fi
}

#x build_install_sh


####################### export the proper environment ##

PATH=$prefix/bin:$prefix/sbin:${PATH}; export PATH
PKG_DBDIR=$pkgdbdir; export PKG_DBDIR
LOCALBASE=$prefix; export LOCALBASE
VARBASE=$varbase; export VARBASE
if [ x"$has_ssp" = x"no" ] && [ x"$check_ssp" = x"yes" ]; then
_OPSYS_SUPPORTS_SSP=no; export _OPSYS_SUPPORTS_SSP
fi


############################## set up an example mk.conf file ####

TARGET_MKCONF=${wrkdir}/mk.conf.example
#x echo_msg "Creating default mk.conf in ${wrkdir}"
#x echo "# Example ${sysconfdir}/mk.conf file produced by bootstrap-pkgsrc" > ${TARGET_MKCONF}
#x echo "# `date`" >> ${TARGET_MKCONF}
#x echo "" >> ${TARGET_MKCONF}
#x echo ".ifdef BSD_PKG_MK	# begin pkgsrc settings" >> ${TARGET_MKCONF}
#x echo "" >> ${TARGET_MKCONF}
#x 
#x if [ -n "$abi" ]; then
#x 	echo "ABI=			$abi" >> ${TARGET_MKCONF}
#x fi
#x if [ "$compiler" != "" ]; then
#x 	echo "PKGSRC_COMPILER=	$compiler" >> ${TARGET_MKCONF}
#x fi
#x case "$compiler" in
#x sunpro)
#x 	echo "CC=			cc"        >> ${TARGET_MKCONF}
#x 	echo "CXX=			CC"        >> ${TARGET_MKCONF}
#x 	echo "CPP=			\${CC} -E" >> ${TARGET_MKCONF}
#x 	;;
#x clang)
#x 	echo "CC=			clang"     >> ${TARGET_MKCONF}
#x 	echo "CXX=			clang++"   >> ${TARGET_MKCONF}
#x 	echo "CPP=			\${CC} -E" >> ${TARGET_MKCONF}
#x 	if [ -z "$CLANGBASE" -a -f "/usr/bin/clang" ]; then
#x 		CLANGBASE="/usr"
#x 	fi
#x 	if [ -n "$CLANGBASE" -o -f "/bin/clang" ]; then
#x 		echo "CLANGBASE=		$CLANGBASE" >> ${TARGET_MKCONF}
#x 	fi
#x 	;;
#x esac
#x if [ -n "$GCCBASE" ]; then
#x 	echo "GCCBASE=		$GCCBASE" >> ${TARGET_MKCONF}
#x fi
#x if [ -n "$SUNWSPROBASE" ]; then
#x 	echo "SUNWSPROBASE=		$SUNWSPROBASE" >> ${TARGET_MKCONF}
#x fi
#x echo "" >> ${TARGET_MKCONF}
#x 
#x if [ x"$has_ssp" = x"no" ] && [ x"$check_ssp" = x"yes" ]; then
#x 	echo "_OPSYS_SUPPORTS_SSP=	no" >> ${TARGET_MKCONF}
#x fi
#x 
#x # for debugging, mainly
#x echo "PKG_DEVELOPER=		yes" >> ${TARGET_MKCONF}
#x echo "ECHO_WRAPPER_MSG=		\${ECHO}" >> ${TARGET_MKCONF}
#x echo "UNPRIVILEGED_USER=		${user}" >> ${TARGET_MKCONF}
#x echo "UNPRIVILEGED_GROUP=		${group}" >> ${TARGET_MKCONF}
#x echo "PKG_DEBUG_LEVEL=		2" >> ${TARGET_MKCONF}
#x 
#x # enable unprivileged builds if not root
#x if [ "$unprivileged" = "yes" ]; then
#x 	echo "UNPRIVILEGED=		yes" >> ${TARGET_MKCONF}
#x fi
#x 
#x # save environment in example mk.conf
#x echo "PKG_DBDIR=		$pkgdbdir" >> ${TARGET_MKCONF}
#x echo "LOCALBASE=		$prefix" >> ${TARGET_MKCONF}
#x if [ "${sysconfbase}" != "/etc" ]; then
#x echo "SYSCONFBASE=		$sysconfbase" >> ${TARGET_MKCONF}
#x fi
#x echo "VARBASE=		$varbase" >> ${TARGET_MKCONF}
#x if [ "${sysconfdir}" != "${prefix}/etc" ]; then
#x 	echo "PKG_SYSCONFBASE=	$sysconfdir" >> ${TARGET_MKCONF}
#x fi
#x echo "PKG_TOOLS_BIN=		$prefix/sbin" >> ${TARGET_MKCONF}
#x echo "PKGINFODIR=		$pkginfodir" >> ${TARGET_MKCONF}
#x echo "PKGMANDIR=		$pkgmandir" >> ${TARGET_MKCONF}
#x echo "" >> ${TARGET_MKCONF}

BOOTSTRAP_MKCONF=${wrkdir}/mk.conf
#x cp ${TARGET_MKCONF} ${BOOTSTRAP_MKCONF}

#x case "$cwrappers" in
#x yes|no)
#x 	echo "USE_CWRAPPERS=		$cwrappers" >> ${TARGET_MKCONF}
#x 	echo "USE_CWRAPPERS=		$cwrappers" >> ${BOOTSTRAP_MKCONF}
#x 	echo "" >> ${TARGET_MKCONF}
#x 	;;
#x esac

# sbin is used by pkg_install, share/mk by bootstrap-mk-files
mkdir -p $wrkdir/sbin $wrkdir/share/mk

###################

bootstrap_mk_files() {
  echo_msg "Bootstrapping mk-files"
  run_cmd "(cd ${pkgsrcdir}/pkgtools/bootstrap-mk-files/files && env CP=cp \
   OPSYS=${opsys} MK_DST=${wrkdir}/share/mk ROOT_GROUP=${root_group} \
  ROOT_USER=${root_user} SED=sed SYSCONFDIR=${sysconfdir} \
  sh ./bootstrap.sh)"
}

#x bootstrap_mk_files

###################

bootstrap_bmake() {
	echo_msg "Bootstrapping bmake"
	copy_src $pkgsrcdir/devel/bmake/files bmake
	run_cmd "(cd $wrkdir/bmake && $shprog configure $configure_quiet_flags --prefix=$wrkdir --with-default-sys-path=$wrkdir/share/mk --with-machine-arch=${machine_arch} $bmakexargs)"
	run_cmd "(cd $wrkdir/bmake && $shprog make-bootstrap.sh)"
	run_cmd "$install_sh -c -o $user -g $group -m 755 $wrkdir/bmake/bmake $wrkdir/bin/bmake"
}

#x bootstrap_bmake

bmake="$wrkdir/bin/bmake"

###################

# build libnbcompat

build_libnbcompat() {
  echo_msg "Building libnbcompat"
  copy_src $pkgsrcdir/pkgtools/libnbcompat/files libnbcompat
  run_cmd "(cd $wrkdir/libnbcompat; $shprog ./configure $configure_quiet_flags -C --prefix=$prefix --infodir=$infodir --mandir=$mandir --sysconfdir=$sysconfdir --enable-bsd-getopt --enable-db && $bmake $make_quiet_flags -j$make_jobs)"
}

#x build_libnbcompat

###################

# bootstrap pkg_install
extra_libarchive_depends() {
	$sedprog -n -e 's/Libs.private: //p' $wrkdir/libarchive/build/pkgconfig/libarchive.pc
}

#x echo_msg "Bootstrapping pkgtools"
#x 
#x copy_src $pkgsrcdir/archivers/libarchive/files libarchive
#x run_cmd "(cd $wrkdir/libarchive; env $BSTRAP_ENV CONFIG_SHELL=$shprog \
#x $shprog ./configure $configure_quiet_flags --enable-static --disable-shared \
#x --disable-bsdcat --disable-bsdtar --disable-bsdcpio --disable-bsdunzip \
#x --disable-posix-regex-lib --disable-xattr --disable-maintainer-mode \
#x --disable-acl --without-zlib --without-bz2lib --without-iconv --without-lzma \
#x --without-lzo2 --without-lz4 --without-nettle --without-openssl \
#x --without-xml2 --without-expat --without-zstd \
#x MAKE=$bmake && $bmake $make_quiet_flags -j$make_jobs)"
#x 
#x copy_src $pkgsrcdir/pkgtools/pkg_install/files pkg_install
#x run_cmd "(cd $wrkdir/pkg_install; env $BSTRAP_ENV \
#x CPPFLAGS='$CPPFLAGS -I${wrkdir}/libnbcompat -I${wrkdir}/libarchive/libarchive' \
#x LDFLAGS='$LDFLAGS -L${wrkdir}/libnbcompat' \
#x LIBS='$LIBS -lnbcompat' $shprog ./configure $configure_quiet_flags -C \
#x --enable-bootstrap --prefix=$prefix --sysconfdir=$sysconfdir \
#x --with-pkgdbdir=$pkgdbdir --infodir=$infodir \
#x --mandir=$mandir $pkg_install_args && \
#x STATIC_LIBARCHIVE=$wrkdir/libarchive/.libs/libarchive.a \
#x STATIC_LIBARCHIVE_LDADD='`extra_libarchive_depends`' \
#x PKGSRC_MACHINE_ARCH="$machine_arch" $bmake $make_quiet_flags -j$make_jobs)"
#x 
#x run_cmd "$install_sh -c -o $user -g $group -m 755 $wrkdir/pkg_install/add/pkg_add $wrkdir/sbin/pkg_add"
#x run_cmd "$install_sh -c -o $user -g $group -m 755 $wrkdir/pkg_install/admin/pkg_admin $wrkdir/sbin/pkg_admin"
#x run_cmd "$install_sh -c -o $user -g $group -m 755 $wrkdir/pkg_install/create/pkg_create $wrkdir/sbin/pkg_create"
#x run_cmd "$install_sh -c -o $user -g $group -m 755 $wrkdir/pkg_install/info/pkg_info $wrkdir/sbin/pkg_info"
#x 
#x echo "NATIVE_PKG_ADD_CMD?=		$wrkdir/sbin/pkg_add" >> ${BOOTSTRAP_MKCONF}
#x echo "NATIVE_PKG_ADMIN_CMD?=		$wrkdir/sbin/pkg_admin" >> ${BOOTSTRAP_MKCONF}
#x echo "NATIVE_PKG_CREATE_CMD?=		$wrkdir/sbin/pkg_create" >> ${BOOTSTRAP_MKCONF}
#x echo "NATIVE_PKG_INFO_CMD?=		$wrkdir/sbin/pkg_info" >> ${BOOTSTRAP_MKCONF}

MAKECONF=$wrkdir/mk.conf
export MAKECONF

###################

#x echo "WRKOBJDIR=		${wrkdir}/wrk" >> ${BOOTSTRAP_MKCONF}
#x 
#x echo "" >> ${TARGET_MKCONF}
#x echo "" >> ${BOOTSTRAP_MKCONF}
#x if test -n "${mk_fragment}"; then
#x 	cat "${mk_fragment}" >> ${TARGET_MKCONF}
#x 	echo "" >> ${TARGET_MKCONF}
#x fi
#x echo ".endif			# end pkgsrc settings" >> ${TARGET_MKCONF}
#x echo ".endif			# end pkgsrc settings" >> ${BOOTSTRAP_MKCONF}

###################

# build and register packages
# usage: build_package <packagedirectory>
build_package() {
	run_cmd "(cd $pkgsrcdir/$1 && $bmake $make_quiet_flags MAKE_JOBS=${make_jobs} PKG_COMPRESSION=none -DPKG_PRESERVE PKGSRC_KEEP_BIN_PKGS=no MAKECONF=${BOOTSTRAP_MKCONF} install)"
}

build_package_nopreserve() {
	run_cmd "(cd $pkgsrcdir/$1 && $bmake $make_quiet_flags MAKE_JOBS=${make_jobs} PKG_COMPRESSION=none PKGSRC_KEEP_BIN_PKGS=no MAKECONF=${BOOTSTRAP_MKCONF} install)"
}

###################

#
# Special packages that we don't want marked with BOOTSTRAP_PKG, but must be
# built (if required) without -DPKG_PRESERVE set so that they can be deleted.
#

use_cwrappers=`(cd $pkgsrcdir/devel/bmake && $bmake show-var VARNAME=_USE_CWRAPPERS)`
case "$use_cwrappers" in
yes)
	build_package_nopreserve "pkgtools/cwrappers"
	;;
esac

use_mktools=`(cd $pkgsrcdir/devel/bmake && $bmake show-var VARNAME=_PKGSRC_USE_MKTOOLS)`
case "$use_mktools" in
yes)
	build_package_nopreserve "pkgtools/mktools"
	;;
esac

#
# Please make sure that the following packages and
# only the following packages set BOOTSTRAP_PKG=yes.
#
echo_msg "Installing packages"
build_package "pkgtools/bootstrap-mk-files"

case "$need_bsd_install" in
yes)
	if [ "$use_bsdinstall" = "yes" ]; then
		build_package "sysutils/bsdinstall"
	else
		build_package "sysutils/install-sh"
	fi
	;;
esac

case "$need_mksh" in
yes)	build_package "shells/mksh";;
esac

case "$need_ksh" in
yes)	build_package "shells/pdksh";;
esac

build_package "devel/bmake"

case "$need_awk" in
yes)	build_package "lang/nawk";;
esac

case "$need_sed" in
yes)	build_package "textproc/nbsed";;
esac

case "$need_extras" in
yes)	build_package "pkgtools/bootstrap-extras";;
esac

case "$need_pax" in
yes)    build_package "archivers/pax"
esac

build_package "pkgtools/pkg_install"

etc_mk_conf="$sysconfdir/mk.conf"

# Install the example mk.conf so that it is used, but only if it doesn't
# exist yet. This can happen with non-default sysconfdir settings.
mkdir_p "$sysconfdir"
if [ ! -f "$etc_mk_conf" ]; then
	cp "$TARGET_MKCONF" "$etc_mk_conf"
	TARGET_MKCONF="$etc_mk_conf"
fi
