# $NetBSD: Makefile,v 1.39 2018/08/16 18:55:02 adam Exp $

DISTNAME=	kmousetool-16.12.0
PKGREVISION=	6
CATEGORIES=	misc
COMMENT=	KDE mouse manipulation tool
MASTER_SITES=	${MASTER_SITE_KDE:=applications/16.12.0/src/}

.include "../../meta-pkgs/kde4/Makefile.kde4"

.include "../../meta-pkgs/kde4/kde4.mk"

.include "../../graphics/hicolor-icon-theme/buildlink3.mk"
.include "../../sysutils/desktop-file-utils/desktopdb.mk"
.include "../../x11/kdelibs4/buildlink3.mk"

.include "../../mk/bsd.pkg.mk"
