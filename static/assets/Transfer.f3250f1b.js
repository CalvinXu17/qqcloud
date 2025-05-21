import{d as L,H as t,dg as Re,dh as ze,di as Te,dj as ye,dk as ke,dl as Fe,dm as Oe,dn as we,dp as Le,d2 as _e,as as H,aS as Pe,p as j,dq as Ve,dr as Ae,ds as $e,n as Be,r as $,dt as Ie,du as Ue,dv as E,_ as He,d4 as De,dw as Ne,bf as U,$ as g,cV as p,dx as A,cW as b,d7 as W,dy as Me,cX as Ee,cY as J,dz as qe,dA as w,dB as je,ar as We,dC as Ke,dc as M}from"./index.59f76cad.js";const Xe=L({name:"Search",render(){return t("svg",{version:"1.1",xmlns:"http://www.w3.org/2000/svg",viewBox:"0 0 512 512",style:"enable-background: new 0 0 512 512"},t("path",{d:`M443.5,420.2L336.7,312.4c20.9-26.2,33.5-59.4,33.5-95.5c0-84.5-68.5-153-153.1-153S64,132.5,64,217s68.5,153,153.1,153
  c36.6,0,70.1-12.8,96.5-34.2l106.1,107.1c3.2,3.4,7.6,5.1,11.9,5.1c4.1,0,8.2-1.5,11.3-4.5C449.5,437.2,449.7,426.8,443.5,420.2z
   M217.1,337.1c-32.1,0-62.3-12.5-85-35.2c-22.7-22.7-35.2-52.9-35.2-84.9c0-32.1,12.5-62.3,35.2-84.9c22.7-22.7,52.9-35.2,85-35.2
  c32.1,0,62.3,12.5,85,35.2c22.7,22.7,35.2,52.9,35.2,84.9c0,32.1-12.5,62.3-35.2,84.9C279.4,324.6,249.2,337.1,217.1,337.1z`}))}}),Ye=e=>{const{fontWeight:a,fontSizeLarge:o,fontSizeMedium:i,fontSizeSmall:d,heightLarge:s,heightMedium:c,borderRadius:l,cardColor:u,tableHeaderColor:f,textColor1:x,textColorDisabled:v,textColor2:k,textColor3:C,borderColor:m,hoverColor:S,closeColorHover:T,closeColorPressed:y,closeIconColor:_,closeIconColorHover:P,closeIconColorPressed:r}=e;return Object.assign(Object.assign({},we),{itemHeightSmall:c,itemHeightMedium:c,itemHeightLarge:s,fontSizeSmall:d,fontSizeMedium:i,fontSizeLarge:o,borderRadius:l,dividerColor:m,borderColor:m,listColor:u,headerColor:Le(u,f),titleTextColor:x,titleTextColorDisabled:v,extraTextColor:C,extraTextColorDisabled:v,itemTextColor:k,itemTextColorDisabled:v,itemColorPending:S,titleFontWeight:a,closeColorHover:T,closeColorPressed:y,closeIconColor:_,closeIconColorHover:P,closeIconColorPressed:r})},Ge=Re({name:"Transfer",common:ze,peers:{Checkbox:Te,Scrollbar:ye,Input:ke,Empty:Fe,Button:Oe},self:Ye}),Je=Ge,B=_e("n-transfer"),K=L({name:"TransferHeader",props:{size:{type:String,required:!0},source:Boolean,onCheckedAll:Function,onClearAll:Function,title:String},setup(e){const{targetOptionsRef:a,canNotSelectAnythingRef:o,canBeClearedRef:i,allCheckedRef:d,mergedThemeRef:s,disabledRef:c,mergedClsPrefixRef:l,srcOptionsLengthRef:u}=H(B),{localeRef:f}=Pe("Transfer");return()=>{const{source:x,onClearAll:v,onCheckedAll:k}=e,{value:C}=s,{value:m}=l,{value:S}=f,T=e.size==="large"?"small":"tiny",{title:y}=e;return t("div",{class:`${m}-transfer-list-header`},y&&t("div",{class:`${m}-transfer-list-header__title`},y),x&&t(j,{class:`${m}-transfer-list-header__button`,theme:C.peers.Button,themeOverrides:C.peerOverrides.Button,size:T,tertiary:!0,onClick:d.value?v:k,disabled:o.value||c.value},{default:()=>d.value?S.unselectAll:S.selectAll}),!x&&i.value&&t(j,{class:`${m}-transfer-list-header__button`,theme:C.peers.Button,themeOverrides:C.peerOverrides.Button,size:T,tertiary:!0,onClick:v,disabled:c.value},{default:()=>S.clearAll}),t("div",{class:`${m}-transfer-list-header__extra`},x?S.total(u.value):S.selected(a.value.length)))}}}),X=L({name:"NTransferListItem",props:{source:Boolean,label:{type:String,required:!0},value:{type:[String,Number],required:!0},disabled:Boolean,option:{type:Object,required:!0}},setup(e){const{targetValueSetRef:a,mergedClsPrefixRef:o,mergedThemeRef:i,handleItemCheck:d,renderSourceLabelRef:s,renderTargetLabelRef:c,showSelectedRef:l}=H(B),u=Ve(()=>a.value.has(e.value));function f(){e.disabled||d(!u.value,e.value)}return{mergedClsPrefix:o,mergedTheme:i,checked:u,showSelected:l,renderSourceLabel:s,renderTargetLabel:c,handleClick:f}},render(){const{disabled:e,mergedTheme:a,mergedClsPrefix:o,label:i,checked:d,source:s,renderSourceLabel:c,renderTargetLabel:l}=this;return t("div",{class:[`${o}-transfer-list-item`,e&&`${o}-transfer-list-item--disabled`,s?`${o}-transfer-list-item--source`:`${o}-transfer-list-item--target`],onClick:s?this.handleClick:void 0},t("div",{class:`${o}-transfer-list-item__background`}),s&&this.showSelected&&t("div",{class:`${o}-transfer-list-item__checkbox`},t(Be,{theme:a.peers.Checkbox,themeOverrides:a.peerOverrides.Checkbox,disabled:e,checked:d})),t("div",{class:`${o}-transfer-list-item__label`,title:$e(i)},s?c?c({option:this.option}):i:l?l({option:this.option}):i),!s&&!e&&t(Ae,{focusable:!1,class:`${o}-transfer-list-item__close`,clsPrefix:o,onClick:this.handleClick}))}}),Y=L({name:"TransferList",props:{virtualScroll:{type:Boolean,required:!0},itemSize:{type:Number,required:!0},options:{type:Array,required:!0},disabled:{type:Boolean,required:!0},source:Boolean},setup(){const{mergedThemeRef:e,mergedClsPrefixRef:a}=H(B),o=$(null),i=$(null);function d(){var l;(l=o.value)===null||l===void 0||l.sync()}function s(){const{value:l}=i;if(!l)return null;const{listElRef:u}=l;return u}function c(){const{value:l}=i;if(!l)return null;const{itemsElRef:u}=l;return u}return{mergedTheme:e,mergedClsPrefix:a,scrollerInstRef:o,vlInstRef:i,syncVLScroller:d,scrollContainer:s,scrollContent:c}},render(){const{mergedTheme:e,options:a}=this;if(a.length===0)return t(Ie,{theme:e.peers.Empty,themeOverrides:e.peerOverrides.Empty});const{mergedClsPrefix:o,virtualScroll:i,source:d,disabled:s,syncVLScroller:c}=this;return t(E,{ref:"scrollerInstRef",theme:e.peers.Scrollbar,themeOverrides:e.peerOverrides.Scrollbar,container:i?this.scrollContainer:void 0,content:i?this.scrollContent:void 0},{default:()=>i?t(Ue,{ref:"vlInstRef",style:{height:"100%"},class:`${o}-transfer-list-content`,items:this.options,itemSize:this.itemSize,showScrollbar:!1,onResize:c,onScroll:c,keyField:"value"},{default:({item:l})=>{const{source:u,disabled:f}=this;return t(X,{source:u,key:l.value,value:l.value,disabled:l.disabled||f,label:l.label,option:l})}}):t("div",{class:`${o}-transfer-list-content`},a.map(l=>t(X,{source:d,key:l.value,value:l.value,disabled:l.disabled||s,label:l.label,option:l})))})}}),G=L({name:"TransferFilter",props:{value:String,placeholder:String,disabled:Boolean,onUpdateValue:{type:Function,required:!0}},setup(){const{mergedThemeRef:e,mergedClsPrefixRef:a}=H(B);return{mergedClsPrefix:a,mergedTheme:e}},render(){const{mergedTheme:e,mergedClsPrefix:a}=this;return t("div",{class:`${a}-transfer-filter`},t(He,{value:this.value,onUpdateValue:this.onUpdateValue,disabled:this.disabled,placeholder:this.placeholder,theme:e.peers.Input,themeOverrides:e.peerOverrides.Input,clearable:!0,size:"small"},{"clear-icon-placeholder":()=>t(De,{clsPrefix:a},{default:()=>t(Xe,null)})}))}});function Qe(e){const a=$(e.defaultValue),o=Ne(U(e,"value"),a),i=g(()=>{const r=new Map;return(e.options||[]).forEach(n=>r.set(n.value,n)),r}),d=g(()=>new Set(o.value||[])),s=g(()=>{const r=i.value,n=[];return(o.value||[]).forEach(V=>{const R=r.get(V);R&&n.push(R)}),n}),c=$(""),l=$(""),u=g(()=>e.sourceFilterable||!!e.filterable),f=g(()=>{const{showSelected:r,options:n,filter:V}=e;return u.value?n.filter(R=>V(c.value,R,"source")&&(r||!d.value.has(R.value))):r?n:n.filter(R=>!d.value.has(R.value))}),x=g(()=>{if(!e.targetFilterable)return s.value;const{filter:r}=e;return s.value.filter(n=>r(l.value,n,"target"))}),v=g(()=>{const{value:r}=o;return r===null?new Set:new Set(r)}),k=g(()=>{const r=new Set(v.value);return f.value.forEach(n=>{!n.disabled&&!r.has(n.value)&&r.add(n.value)}),r}),C=g(()=>{const r=new Set(v.value);return f.value.forEach(n=>{!n.disabled&&r.has(n.value)&&r.delete(n.value)}),r}),m=g(()=>{const r=new Set(v.value);return x.value.forEach(n=>{n.disabled||r.delete(n.value)}),r}),S=g(()=>f.value.every(r=>r.disabled)),T=g(()=>{if(!f.value.length)return!1;const r=v.value;return f.value.every(n=>n.disabled||r.has(n.value))}),y=g(()=>x.value.some(r=>!r.disabled));function _(r){c.value=r!=null?r:""}function P(r){l.value=r!=null?r:""}return{uncontrolledValueRef:a,mergedValueRef:o,targetValueSetRef:d,valueSetForCheckAllRef:k,valueSetForUncheckAllRef:C,valueSetForClearRef:m,filteredTgtOptionsRef:x,filteredSrcOptionsRef:f,targetOptionsRef:s,canNotSelectAnythingRef:S,canBeClearedRef:y,allCheckedRef:T,srcPatternRef:c,tgtPatternRef:l,mergedSrcFilterableRef:u,handleSrcFilterUpdateValue:_,handleTgtFilterUpdateValue:P}}const Ze=p("transfer",`
 width: 100%;
 font-size: var(--n-font-size);
 height: 300px;
 display: flex;
 flex-wrap: nowrap;
 word-break: break-word;
`,[A("disabled",[p("transfer-list",[p("transfer-list-header",[b("title",`
 color: var(--n-header-text-color-disabled);
 `),b("extra",`
 color: var(--n-header-extra-text-color-disabled);
 `)])])]),p("transfer-list",`
 flex: 1;
 min-width: 0;
 height: inherit;
 display: flex;
 flex-direction: column;
 background-clip: padding-box;
 position: relative;
 transition: background-color .3s var(--n-bezier);
 background-color: var(--n-list-color);
 `,[A("source",`
 border-top-left-radius: var(--n-border-radius);
 border-bottom-left-radius: var(--n-border-radius);
 `,[b("border","border-right: 1px solid var(--n-divider-color);")]),A("target",`
 border-top-right-radius: var(--n-border-radius);
 border-bottom-right-radius: var(--n-border-radius);
 `,[b("border","border-left: none;")]),b("border",`
 padding: 0 12px;
 border: 1px solid var(--n-border-color);
 transition: border-color .3s var(--n-bezier);
 pointer-events: none;
 border-radius: inherit;
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 `),p("transfer-list-header",`
 min-height: var(--n-header-height);
 box-sizing: border-box;
 display: flex;
 padding: 12px 12px 10px 12px;
 align-items: center;
 background-clip: padding-box;
 border-radius: inherit;
 border-bottom-left-radius: 0;
 border-bottom-right-radius: 0;
 line-height: 1.5;
 transition:
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `,[W("> *:not(:first-child)",`
 margin-left: 8px;
 `),b("title",`
 flex: 1;
 min-width: 0;
 line-height: 1.5;
 font-size: var(--n-header-font-size);
 font-weight: var(--n-header-font-weight);
 transition: color .3s var(--n-bezier);
 color: var(--n-header-text-color);
 `),b("button",`
 position: relative;
 `),b("extra",`
 transition: color .3s var(--n-bezier);
 font-size: var(--n-extra-font-size);
 margin-right: 0;
 white-space: nowrap;
 color: var(--n-header-extra-text-color);
 `)]),p("transfer-list-body",`
 flex-basis: 0;
 flex-grow: 1;
 box-sizing: border-box;
 position: relative;
 display: flex;
 flex-direction: column;
 border-radius: inherit;
 border-top-left-radius: 0;
 border-top-right-radius: 0;
 `,[p("transfer-filter",`
 padding: 4px 12px 8px 12px;
 box-sizing: border-box;
 transition:
 border-color .3s var(--n-bezier),
 background-color .3s var(--n-bezier);
 `),p("transfer-list-flex-container",`
 flex: 1;
 position: relative;
 `,[p("scrollbar",`
 position: absolute;
 left: 0;
 right: 0;
 top: 0;
 bottom: 0;
 height: unset;
 `),p("empty",`
 position: absolute;
 left: 50%;
 top: 50%;
 transform: translateY(-50%) translateX(-50%);
 `),p("transfer-list-content",`
 padding: 0;
 margin: 0;
 position: relative;
 `,[p("transfer-list-item",`
 padding: 0 12px;
 min-height: var(--n-item-height);
 display: flex;
 align-items: center;
 color: var(--n-item-text-color);
 position: relative;
 transition: color .3s var(--n-bezier);
 `,[b("background",`
 position: absolute;
 left: 4px;
 right: 4px;
 top: 0;
 bottom: 0;
 border-radius: var(--n-border-radius);
 transition: background-color .3s var(--n-bezier);
 `),b("checkbox",`
 position: relative;
 margin-right: 8px;
 `),b("close",`
 opacity: 0;
 pointer-events: none;
 position: relative;
 transition:
 opacity .3s var(--n-bezier),
 background-color .3s var(--n-bezier),
 color .3s var(--n-bezier);
 `),b("label",`
 position: relative;
 min-width: 0;
 flex-grow: 1;
 `),A("source","cursor: pointer;"),A("disabled",`
 cursor: not-allowed;
 color: var(--n-item-text-color-disabled);
 `),Me("disabled",[W("&:hover",[b("background","background-color: var(--n-item-color-pending);"),b("close",`
 opacity: 1;
 pointer-events: all;
 `)])])])])])])])]),er=Object.assign(Object.assign({},J.props),{value:Array,defaultValue:{type:Array,default:null},options:{type:Array,default:()=>[]},disabled:{type:Boolean,default:void 0},virtualScroll:Boolean,sourceTitle:String,targetTitle:String,filterable:{type:Boolean,default:void 0},sourceFilterable:Boolean,targetFilterable:Boolean,showSelected:{type:Boolean,default:!0},sourceFilterPlaceholder:String,targetFilterPlaceholder:String,filter:{type:Function,default:(e,a)=>e?~(""+a.label).toLowerCase().indexOf((""+e).toLowerCase()):!0},size:String,renderSourceLabel:Function,renderTargetLabel:Function,renderSourceList:Function,renderTargetList:Function,"onUpdate:value":[Function,Array],onUpdateValue:[Function,Array],onChange:[Function,Array]}),tr=L({name:"Transfer",props:er,setup(e){const{mergedClsPrefixRef:a}=Ee(e),o=J("Transfer","-transfer",Ze,Je,e,a),i=qe(e),{mergedSizeRef:d,mergedDisabledRef:s}=i,c=g(()=>{const{value:h}=d,{self:{[w("itemHeight",h)]:z}}=o.value;return je(z)}),{uncontrolledValueRef:l,mergedValueRef:u,targetValueSetRef:f,valueSetForCheckAllRef:x,valueSetForUncheckAllRef:v,valueSetForClearRef:k,filteredTgtOptionsRef:C,filteredSrcOptionsRef:m,targetOptionsRef:S,canNotSelectAnythingRef:T,canBeClearedRef:y,allCheckedRef:_,srcPatternRef:P,tgtPatternRef:r,mergedSrcFilterableRef:n,handleSrcFilterUpdateValue:V,handleTgtFilterUpdateValue:R}=Qe(e);function F(h){const{onUpdateValue:z,"onUpdate:value":O,onChange:I}=e,{nTriggerFormInput:D,nTriggerFormChange:N}=i;z&&M(z,h),O&&M(O,h),I&&M(I,h),l.value=h,D(),N()}function Q(){F([...x.value])}function Z(){F([...v.value])}function ee(){F([...k.value])}function q(h,z){F(h?(u.value||[]).concat(z):(u.value||[]).filter(O=>O!==z))}function re(h){F(h)}return We(B,{targetValueSetRef:f,mergedClsPrefixRef:a,disabledRef:s,mergedThemeRef:o,targetOptionsRef:S,canNotSelectAnythingRef:T,canBeClearedRef:y,allCheckedRef:_,srcOptionsLengthRef:g(()=>e.options.length),handleItemCheck:q,renderSourceLabelRef:U(e,"renderSourceLabel"),renderTargetLabelRef:U(e,"renderTargetLabel"),showSelectedRef:U(e,"showSelected")}),{mergedClsPrefix:a,mergedDisabled:s,itemSize:c,isMounted:Ke(),mergedTheme:o,filteredSrcOpts:m,filteredTgtOpts:C,srcPattern:P,tgtPattern:r,mergedSize:d,mergedSrcFilterable:n,handleSrcFilterUpdateValue:V,handleTgtFilterUpdateValue:R,handleSourceCheckAll:Q,handleSourceUncheckAll:Z,handleTargetClearAll:ee,handleItemCheck:q,handleChecked:re,cssVars:g(()=>{const{value:h}=d,{common:{cubicBezierEaseInOut:z},self:{borderRadius:O,borderColor:I,listColor:D,titleTextColor:N,titleTextColorDisabled:te,extraTextColor:le,itemTextColor:oe,itemColorPending:ae,itemTextColorDisabled:ie,titleFontWeight:ne,closeColorHover:se,closeColorPressed:de,closeIconColor:ce,closeIconColorHover:ue,closeIconColorPressed:fe,closeIconSize:he,closeSize:ge,dividerColor:be,extraTextColorDisabled:ve,[w("extraFontSize",h)]:me,[w("fontSize",h)]:pe,[w("titleFontSize",h)]:xe,[w("itemHeight",h)]:Se,[w("headerHeight",h)]:Ce}}=o.value;return{"--n-bezier":z,"--n-border-color":I,"--n-border-radius":O,"--n-extra-font-size":me,"--n-font-size":pe,"--n-header-font-size":xe,"--n-header-extra-text-color":le,"--n-header-extra-text-color-disabled":ve,"--n-header-font-weight":ne,"--n-header-text-color":N,"--n-header-text-color-disabled":te,"--n-item-color-pending":ae,"--n-item-height":Se,"--n-item-text-color":oe,"--n-item-text-color-disabled":ie,"--n-list-color":D,"--n-header-height":Ce,"--n-close-size":ge,"--n-close-icon-size":he,"--n-close-color-hover":se,"--n-close-color-pressed":de,"--n-close-icon-color":ce,"--n-close-icon-color-hover":ue,"--n-close-icon-color-pressed":fe,"--n-divider-color":be}})}},render(){const{mergedClsPrefix:e,renderSourceList:a,renderTargetList:o,mergedTheme:i,mergedSrcFilterable:d,targetFilterable:s}=this;return t("div",{class:[`${e}-transfer`,this.mergedDisabled&&`${e}-transfer--disabled`],style:this.cssVars},t("div",{class:`${e}-transfer-list ${e}-transfer-list--source`},t(K,{source:!0,title:this.sourceTitle,onCheckedAll:this.handleSourceCheckAll,onClearAll:this.handleSourceUncheckAll,size:this.mergedSize}),t("div",{class:`${e}-transfer-list-body`},d?t(G,{onUpdateValue:this.handleSrcFilterUpdateValue,value:this.srcPattern,disabled:this.mergedDisabled,placeholder:this.sourceFilterPlaceholder}):null,t("div",{class:`${e}-transfer-list-flex-container`},a?t(E,{theme:i.peers.Scrollbar,themeOverrides:i.peerOverrides.Scrollbar},{default:()=>a({onCheck:this.handleChecked,checkedOptions:this.filteredTgtOpts,pattern:this.srcPattern})}):t(Y,{source:!0,options:this.filteredSrcOpts,disabled:this.mergedDisabled,virtualScroll:this.virtualScroll,itemSize:this.itemSize}))),t("div",{class:`${e}-transfer-list__border`})),t("div",{class:`${e}-transfer-list ${e}-transfer-list--target`},t(K,{onClearAll:this.handleTargetClearAll,size:this.mergedSize,title:this.targetTitle}),t("div",{class:`${e}-transfer-list-body`},s?t(G,{onUpdateValue:this.handleTgtFilterUpdateValue,value:this.tgtPattern,disabled:this.mergedDisabled,placeholder:this.sourceFilterPlaceholder}):null,t("div",{class:`${e}-transfer-list-flex-container`},o?t(E,{theme:i.peers.Scrollbar,themeOverrides:i.peerOverrides.Scrollbar},{default:()=>o({onCheck:this.handleChecked,checkedOptions:this.filteredTgtOpts,pattern:this.tgtPattern})}):t(Y,{options:this.filteredTgtOpts,disabled:this.mergedDisabled,virtualScroll:this.virtualScroll,itemSize:this.itemSize}))),t("div",{class:`${e}-transfer-list__border`})))}});export{tr as _};
