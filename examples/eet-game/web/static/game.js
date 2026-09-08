// ============================================================
// EET Space Invaders – browser client
// Connects via WebSocket, renders game state on HTML5 Canvas
// ============================================================
'use strict';

const W=800,H=600,PLAYER_W=40,PLAYER_H=24,BULLET_W=4,BULLET_H=12,
      SHIELD_W=64,SHIELD_H=40,UFO_W=48,UFO_H=22;

const C={bg:'#050510',cyan:'#00ffff',magenta:'#ff00ff',green:'#39ff14',
         yellow:'#ffff00',red:'#ff2244',orange:'#ff8800',white:'#ffffff'};

const canvas=document.getElementById('game-canvas');
const ctx=canvas.getContext('2d');
canvas.width=W; canvas.height=H;

let gameState=null,ws=null,shakeAmt=0;
const particles=[];
const stars=Array.from({length:120},()=>({
  x:Math.random()*W, y:Math.random()*H,
  r:Math.random()*1.4+0.3, speed:Math.random()*15+5, bright:Math.random()
}));
const keys={}, held={};

// ── WebSocket ────────────────────────────────────────────────
function connect(){
  const proto=location.protocol==='https:'?'wss':'ws';
  ws=new WebSocket(`${proto}://${location.host}/ws`);
  ws.addEventListener('open',()=>console.log('🟢 WS connected'));
  ws.addEventListener('message',(e)=>{
    const msg=JSON.parse(e.data);
    if(msg.type==='state'){gameState=msg;updateHUD(msg);}
    else if(msg.type==='event') handleEvent(msg);
  });
  ws.addEventListener('close',()=>setTimeout(connect,1500));
  ws.addEventListener('error',(e)=>console.error('ws',e));
}
function send(action){
  if(ws&&ws.readyState===WebSocket.OPEN)
    ws.send(JSON.stringify({type:'input',action}));
}

// ── Input ────────────────────────────────────────────────────
document.addEventListener('keydown',(e)=>{
  if(keys[e.code])return; keys[e.code]=true;
  if(e.code==='ArrowLeft' ||e.code==='KeyA'){startHeld('move_left');}
  else if(e.code==='ArrowRight'||e.code==='KeyD'){startHeld('move_right');}
  else if(e.code==='Space'){e.preventDefault();send('shoot');}
  else if(e.code==='KeyP') send('pause');
  else if(e.code==='Enter'){send('start');send('restart');}
});
document.addEventListener('keyup',(e)=>{
  keys[e.code]=false;
  if(e.code==='ArrowLeft' ||e.code==='KeyA') stopHeld('move_left');
  if(e.code==='ArrowRight'||e.code==='KeyD') stopHeld('move_right');
});
function startHeld(a){if(held[a])return;send(a);held[a]=setInterval(()=>send(a),16);}
function stopHeld(a){clearInterval(held[a]);delete held[a];}

document.getElementById('btn-left').addEventListener('touchstart', e=>{e.preventDefault();startHeld('move_left');});
document.getElementById('btn-left').addEventListener('touchend',   e=>{e.preventDefault();stopHeld('move_left');});
document.getElementById('btn-right').addEventListener('touchstart',e=>{e.preventDefault();startHeld('move_right');});
document.getElementById('btn-right').addEventListener('touchend',  e=>{e.preventDefault();stopHeld('move_right');});
document.getElementById('btn-fire').addEventListener('touchstart', e=>{e.preventDefault();send('shoot');});
document.getElementById('btn-start').addEventListener('touchstart',e=>{e.preventDefault();send('start');send('restart');});

// ── HUD ──────────────────────────────────────────────────────
const elScore=document.getElementById('score-val');
const elHi=document.getElementById('hi-val');
const elLevel=document.getElementById('level-val');
const elLives=document.getElementById('lives-val');
function updateHUD(s){
  elScore.textContent=String(s.score).padStart(6,'0');
  elHi.textContent=String(s.hiScore).padStart(6,'0');
  elLevel.textContent=s.level||'—';
  elLives.textContent='🚀'.repeat(Math.max(0,s.player?.lives??0));
}

// ── Events ───────────────────────────────────────────────────
function handleEvent(msg){
  switch(msg.event){
    case'alien_killed':
      spawnExplosion(msg.data.x,msg.data.y,C.cyan,12);
      showScorePopup(msg.data.x,msg.data.y,'+'+msg.data.score,C.cyan);
      playTone(440,0.08,'square'); break;
    case'ufo_killed':
      spawnExplosion(msg.data.x,msg.data.y,C.magenta,20);
      showScorePopup(msg.data.x,msg.data.y,'+'+msg.data.score,C.yellow);
      playTone(880,0.15,'sawtooth'); break;
    case'player_hit':
      shakeAmt=14;
      spawnExplosion(msg.data.x,msg.data.y,C.red,22);
      playTone(110,0.3,'sawtooth'); break;
    case'shield_hit':
      spawnExplosion(msg.data.x,msg.data.y,C.orange,6);
      playTone(200,0.05,'square'); break;
    case'level_complete':
      showLevelBanner('WAVE '+msg.data.score+' CLEAR!');
      playTone(660,0.2,'sine'); break;
    case'game_over': playTone(55,0.5,'sawtooth'); break;
    case'victory':   playVictory(); break;
  }
}

// ── Audio ────────────────────────────────────────────────────
const audioCtx=new(window.AudioContext||window.webkitAudioContext)();
function playTone(freq,dur,type='square',vol=0.18){
  try{
    const osc=audioCtx.createOscillator(),gain=audioCtx.createGain();
    osc.type=type; osc.frequency.value=freq;
    gain.gain.setValueAtTime(vol,audioCtx.currentTime);
    gain.gain.exponentialRampToValueAtTime(0.001,audioCtx.currentTime+dur);
    osc.connect(gain); gain.connect(audioCtx.destination);
    osc.start(); osc.stop(audioCtx.currentTime+dur);
  }catch(_){}
}
function playVictory(){
  [523,659,784,1047].forEach((f,i)=>setTimeout(()=>playTone(f,0.3,'sine',0.25),i*200));
}

// ── Particles ────────────────────────────────────────────────
function spawnExplosion(x,y,color,count){
  for(let i=0;i<count;i++){
    const a=Math.random()*Math.PI*2,sp=Math.random()*120+40;
    particles.push({x,y,vx:Math.cos(a)*sp,vy:Math.sin(a)*sp,
      life:1,decay:Math.random()*1.5+1,r:Math.random()*3+1,color});
  }
}
function updateParticles(dt){
  for(let i=particles.length-1;i>=0;i--){
    const p=particles[i];
    p.x+=p.vx*dt; p.y+=p.vy*dt; p.vy+=60*dt; p.life-=p.decay*dt;
    if(p.life<=0) particles.splice(i,1);
  }
}
function drawParticles(){
  for(const p of particles){
    ctx.globalAlpha=Math.max(0,p.life);
    ctx.fillStyle=p.color; ctx.shadowColor=p.color; ctx.shadowBlur=6;
    ctx.beginPath(); ctx.arc(p.x,p.y,p.r,0,Math.PI*2); ctx.fill();
  }
  ctx.globalAlpha=1; ctx.shadowBlur=0;
}

// ── Score popups ─────────────────────────────────────────────
const wrapper=document.getElementById('game-wrapper');
function showScorePopup(x,y,text,color){
  const el=document.createElement('div');
  el.className='score-popup'; el.textContent=text; el.style.color=color;
  el.style.left=(x*(wrapper.offsetWidth/W)-20)+'px';
  el.style.top=(y*(wrapper.offsetHeight/H)-10)+'px';
  wrapper.appendChild(el);
  el.addEventListener('animationend',()=>el.remove());
}
function showLevelBanner(text){
  const old=document.getElementById('level-banner');
  if(old)old.remove();
  const el=document.createElement('div');
  el.id='level-banner'; el.textContent=text;
  wrapper.appendChild(el);
  el.addEventListener('animationend',()=>el.remove());
}

// ── Stars ────────────────────────────────────────────────────
function drawStars(dt){
  for(const s of stars){
    s.y+=s.speed*dt; if(s.y>H){s.y=0;s.x=Math.random()*W;}
    ctx.globalAlpha=0.4+0.6*s.bright*Math.abs(Math.sin(Date.now()/1200+s.x));
    ctx.fillStyle=C.white;
    ctx.beginPath(); ctx.arc(s.x,s.y,s.r,0,Math.PI*2); ctx.fill();
  }
  ctx.globalAlpha=1;
}

// ── Helpers ──────────────────────────────────────────────────
function glow(color,blur=12){ctx.shadowColor=color;ctx.shadowBlur=blur;}
function noGlow(){ctx.shadowBlur=0;}

// ── Player ───────────────────────────────────────────────────
function drawPlayer(p){
  if(!p)return;
  if(p.invincible&&Math.floor(Date.now()/100)%2===0)return;
  const x=p.x,y=p.y;
  ctx.save();
  glow(C.cyan,18); ctx.fillStyle=C.cyan;
  ctx.beginPath();
  ctx.moveTo(x+PLAYER_W/2,y); ctx.lineTo(x+PLAYER_W,y+PLAYER_H); ctx.lineTo(x,y+PLAYER_H);
  ctx.closePath(); ctx.fill();
  glow(C.magenta,10); ctx.fillStyle=C.magenta;
  ctx.beginPath(); ctx.arc(x+PLAYER_W/2,y+PLAYER_H*0.55,5,0,Math.PI*2); ctx.fill();
  const fh=8+4*Math.sin(Date.now()/60);
  const grad=ctx.createLinearGradient(x+PLAYER_W/2,y+PLAYER_H,x+PLAYER_W/2,y+PLAYER_H+fh);
  grad.addColorStop(0,C.orange); grad.addColorStop(1,'rgba(255,136,0,0)');
  ctx.fillStyle=grad; glow(C.orange,14);
  ctx.beginPath();
  ctx.moveTo(x+PLAYER_W/2-6,y+PLAYER_H);
  ctx.lineTo(x+PLAYER_W/2+6,y+PLAYER_H);
  ctx.lineTo(x+PLAYER_W/2,y+PLAYER_H+fh);
  ctx.closePath(); ctx.fill();
  ctx.restore();
}

// ── Aliens ───────────────────────────────────────────────────
function drawAlienBottom(x,y,f,c){
  ctx.fillStyle=c;
  if(f===0){
    ctx.fillRect(x+8,y+4,20,12);ctx.fillRect(x+4,y+8,28,6);
    ctx.fillStyle=C.bg;ctx.fillRect(x+10,y+6,4,4);ctx.fillRect(x+22,y+6,4,4);
    ctx.fillStyle=c;
    ctx.fillRect(x+6,y+16,4,6);ctx.fillRect(x+14,y+16,4,6);ctx.fillRect(x+22,y+16,4,6);
    ctx.fillRect(x+10,y,4,4);ctx.fillRect(x+22,y,4,4);
  }else{
    ctx.fillRect(x+8,y+4,20,12);ctx.fillRect(x+4,y+8,28,6);
    ctx.fillStyle=C.bg;ctx.fillRect(x+10,y+6,4,4);ctx.fillRect(x+22,y+6,4,4);
    ctx.fillStyle=c;
    ctx.fillRect(x+4,y+16,4,6);ctx.fillRect(x+16,y+16,4,6);ctx.fillRect(x+28,y+16,4,6);
    ctx.fillRect(x+8,y,4,4);ctx.fillRect(x+24,y,4,4);
  }
}
function drawAlienMiddle(x,y,f,c){
  ctx.fillStyle=c;
  if(f===0){
    ctx.fillRect(x+6,y+2,24,14);ctx.fillRect(x+2,y+8,32,6);
    ctx.fillStyle=C.bg;ctx.fillRect(x+8,y+4,4,4);ctx.fillRect(x+24,y+4,4,4);
    ctx.fillStyle=c;ctx.fillRect(x+2,y+18,6,4);ctx.fillRect(x+28,y+18,6,4);ctx.fillRect(x+10,y,16,4);
  }else{
    ctx.fillRect(x+6,y+2,24,14);ctx.fillRect(x+2,y+8,32,6);
    ctx.fillStyle=C.bg;ctx.fillRect(x+8,y+4,4,4);ctx.fillRect(x+24,y+4,4,4);
    ctx.fillStyle=c;ctx.fillRect(x+0,y+18,6,4);ctx.fillRect(x+30,y+18,6,4);ctx.fillRect(x+12,y,12,4);
  }
}
function drawAlienTop(x,y,f,c){
  ctx.fillStyle=c;
  if(f===0){
    ctx.fillRect(x+12,y,12,4);ctx.fillRect(x+8,y+4,20,12);ctx.fillRect(x+4,y+8,28,4);
    ctx.fillStyle=C.bg;ctx.fillRect(x+10,y+6,4,4);ctx.fillRect(x+22,y+6,4,4);
    ctx.fillStyle=c;ctx.fillRect(x+4,y+16,4,6);ctx.fillRect(x+28,y+16,4,6);
  }else{
    ctx.fillRect(x+12,y,12,4);ctx.fillRect(x+8,y+4,20,12);ctx.fillRect(x+4,y+8,28,4);
    ctx.fillStyle=C.bg;ctx.fillRect(x+10,y+6,4,4);ctx.fillRect(x+22,y+6,4,4);
    ctx.fillStyle=c;ctx.fillRect(x+8,y+18,4,4);ctx.fillRect(x+24,y+18,4,4);
  }
}
const ALIEN_FN={1:drawAlienBottom,2:drawAlienMiddle,3:drawAlienTop};
const ALIEN_COL={1:C.green,2:C.cyan,3:C.magenta};
function drawAliens(aliens){
  if(!aliens)return;
  for(const a of aliens){
    if(!a.alive)continue;
    const col=ALIEN_COL[a.alienType]||C.cyan;
    ctx.save();glow(col,10);(ALIEN_FN[a.alienType]||drawAlienBottom)(a.x,a.y,a.frame,col);ctx.restore();
  }
}

// ── UFO ──────────────────────────────────────────────────────
function drawUFO(ufo){
  if(!ufo||!ufo.active)return;
  const x=ufo.x,y=ufo.y;
  ctx.save();glow(C.red,18);
  const g=ctx.createRadialGradient(x+UFO_W/2,y+UFO_H/2,2,x+UFO_W/2,y+UFO_H/2,UFO_W/2);
  g.addColorStop(0,C.magenta);g.addColorStop(1,C.red);ctx.fillStyle=g;
  ctx.beginPath();ctx.ellipse(x+UFO_W/2,y+UFO_H*.65,UFO_W/2,UFO_H*.4,0,0,Math.PI*2);ctx.fill();
  glow(C.cyan,8);ctx.fillStyle=C.cyan;
  ctx.beginPath();ctx.ellipse(x+UFO_W/2,y+UFO_H*.45,UFO_W*.32,UFO_H*.38,0,Math.PI,0);ctx.fill();
  const t=Date.now()/200;
  [-14,-6,2,10].forEach((ox,i)=>{
    ctx.fillStyle=(Math.floor(t+i)%3===0)?C.yellow:C.red;glow(ctx.fillStyle,6);
    ctx.beginPath();ctx.arc(x+UFO_W/2+ox,y+UFO_H*.7,2.5,0,Math.PI*2);ctx.fill();
  });
  ctx.restore();
}

// ── Bullets ──────────────────────────────────────────────────
function drawBullets(bullets){
  if(!bullets)return;
  for(const b of bullets){
    ctx.save();
    if(b.owner==='player'){
      glow(C.green,14);ctx.fillStyle=C.green;ctx.fillRect(b.x,b.y,BULLET_W,BULLET_H);
      ctx.fillStyle=C.white;ctx.fillRect(b.x+1,b.y,2,3);
    }else{
      glow(C.red,10);ctx.fillStyle=C.red;ctx.fillRect(b.x,b.y,BULLET_W,BULLET_H);
      ctx.fillStyle=C.orange;ctx.fillRect(b.x+1,b.y+3,2,3);ctx.fillRect(b.x+1,b.y+8,2,3);
    }
    ctx.restore();
  }
}

// ── Shields ──────────────────────────────────────────────────
function drawShields(shields){
  if(!shields)return;
  for(const s of shields){
    if(s.hp<=0)continue;
    ctx.save();ctx.globalAlpha=0.3+(s.hp/4)*.7;glow(C.green,8);ctx.fillStyle=C.green;
    ctx.fillRect(s.x,s.y+10,SHIELD_W,SHIELD_H-10);ctx.fillRect(s.x+10,s.y,SHIELD_W-20,10);
    ctx.fillStyle=C.bg;ctx.fillRect(s.x+20,s.y+SHIELD_H-16,24,16);
    if(s.hp<4){
      ctx.strokeStyle=C.bg;ctx.lineWidth=2;
      for(let i=0;i<4-s.hp;i++){
        const cx=s.x+10+i*14+Math.sin(i*3)*4;
        ctx.beginPath();ctx.moveTo(cx,s.y+8);ctx.lineTo(cx+4,s.y+20);ctx.lineTo(cx+2,s.y+32);ctx.stroke();
      }
    }
    ctx.restore();
  }
}

function drawGround(){glow(C.green,4);ctx.fillStyle=C.green;ctx.fillRect(0,H-26,W,2);noGlow();}

// ── Overlays ─────────────────────────────────────────────────
const overlayMenu    =document.getElementById('overlay-menu');
const overlayPaused  =document.getElementById('overlay-paused');
const overlayGameOver=document.getElementById('overlay-gameover');
const overlayVictory =document.getElementById('overlay-victory');
const elGoScore      =document.getElementById('go-score');
const elVicScore     =document.getElementById('vic-score');
let lastOvState='';
function updateOverlays(state,score){
  if(state===lastOvState)return; lastOvState=state;
  overlayMenu.classList.toggle('hidden',    state!=='menu');
  overlayPaused.classList.toggle('hidden',  state!=='paused');
  overlayGameOver.classList.toggle('hidden',state!=='game_over');
  overlayVictory.classList.toggle('hidden', state!=='victory');
  if(state==='game_over'&&elGoScore)  elGoScore.textContent ='SCORE  '+String(score).padStart(6,'0');
  if(state==='victory'  &&elVicScore) elVicScore.textContent='SCORE  '+String(score).padStart(6,'0');
}

// ── Main render loop ─────────────────────────────────────────
let lastTime=0;
function frame(ts){
  const dt=Math.min((ts-lastTime)/1000,.05); lastTime=ts;
  let sx=0,sy=0;
  if(shakeAmt>0){
    sx=(Math.random()-.5)*shakeAmt; sy=(Math.random()-.5)*shakeAmt;
    shakeAmt*=.82; if(shakeAmt<.5)shakeAmt=0;
  }
  ctx.save(); ctx.translate(sx,sy);
  ctx.fillStyle=C.bg; ctx.fillRect(-10,-10,W+20,H+20);
  drawStars(dt); updateParticles(dt);
  if(gameState){
    const s=gameState;
    updateOverlays(s.state,s.score);
    if(s.state==='playing'||s.state==='paused'){
      drawGround(); drawShields(s.shields); drawAliens(s.aliens);
      drawUFO(s.ufo); drawBullets(s.bullets); drawPlayer(s.player);
    }
  }else{updateOverlays('menu',0);}
  drawParticles(); ctx.restore();
  requestAnimationFrame(frame);
}

connect();
requestAnimationFrame(frame);
