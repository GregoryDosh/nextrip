package main

import (
	"github.com/valyala/fasthttp"
)

var graphiqlContent = []byte(`<!DOCTYPE html><title>GraphiQL</title><style>body{height:100%;margin:0;width:100%;overflow:hidden}#graphiql{height:100vh}</style><script src=//cdn.jsdelivr.net/es6-promise/4.0.5/es6-promise.auto.min.js></script><script src=//cdn.jsdelivr.net/fetch/0.9.0/fetch.min.js></script><script src=//cdn.jsdelivr.net/react/15.4.2/react.min.js></script><script src=//cdn.jsdelivr.net/react/15.4.2/react-dom.min.js></script><link href=//cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.css rel=stylesheet><script src=//cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.11/graphiql.min.js></script><div id=graphiql>Loading...</div><script>var search=window.location.search,parameters={};if(search.substr(1).split("&").forEach(function(e){var a=e.indexOf("=");0<=a&&(parameters[decodeURIComponent(e.slice(0,a))]=decodeURIComponent(e.slice(a+1)))}),parameters.variables)try{parameters.variables=JSON.stringify(JSON.parse(parameters.variables),null,2)}catch(e){}function onEditQuery(e){parameters.query=e,updateURL()}function onEditVariables(e){parameters.variables=e,updateURL()}function onEditOperationName(e){parameters.operationName=e,updateURL()}function updateURL(){var e="?"+Object.keys(parameters).filter(function(e){return Boolean(parameters[e])}).map(function(e){return encodeURIComponent(e)+"="+encodeURIComponent(parameters[e])}).join("&");history.replaceState(null,null,e)}function graphQLFetcher(e){return fetch("/graphql",{method:"post",headers:{Accept:"application/json","Content-Type":"application/json"},body:JSON.stringify(e),credentials:"include"}).then(function(e){return e.text()}).then(function(a){try{return JSON.parse(a)}catch(e){return a}})}ReactDOM.render(React.createElement(GraphiQL,{fetcher:graphQLFetcher,query:parameters.query,variables:parameters.variables,operationName:parameters.operationName,onEditQuery:onEditQuery,onEditVariables:onEditVariables,onEditOperationName:onEditOperationName}),document.getElementById("graphiql"))</script>`)

func graphiqlHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	ctx.SetBody(graphiqlContent)
}
