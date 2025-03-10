package server

import (
	"fmt"
	"go/types"
	"path"

	"github.com/goplus/gop/ast"
	gopast "github.com/goplus/gop/ast"
	"github.com/goplus/gop/x/typesutil"
	"github.com/goplus/goxlsw/internal/vfs"
	"github.com/goplus/goxlsw/pkgdoc"
)

type Context struct {
	Proj *vfs.MapFS
	URI  DocumentURI
	File string
	Pos  Position // optional

	Ast  *ast.File
	Info *typesutil.Info
	Doc  *pkgdoc.PkgDoc
}

// Context returns the context for the given document URI and position.
func (s *Server) Context(uri DocumentURI, pos Position) (ctx *Context, err error) {
	spxFile, err := s.fromDocumentURI(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to get file path from document URI %q: %w", uri, err)
	}
	if path.Ext(spxFile) != ".spx" {
		return nil, fmt.Errorf("file %q does not have .spx extension", spxFile)
	}

	proj := s.getProj()
	return &Context{
		Proj: proj,
		URI:  uri,
		Pos:  pos,
		File: spxFile,
		Ast:  getASTPkg(proj).Files[spxFile],
		Info: getTypeInfo(proj),
		Doc:  getPkgDoc(proj),
	}, nil
}

// defIdentFor returns the identifier where the given object is defined.
func (s *Server) defIdentFor(ctx *Context, obj types.Object) *gopast.Ident {
	if obj == nil {
		return nil
	}
	for ident, o := range getTypeInfo(ctx.Proj).Defs {
		if o == obj {
			return ident
		}
	}
	return nil
}
