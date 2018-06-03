import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ContributeAmountComponent } from './contribute-amount.component';

describe('ContributeAmountComponent', () => {
  let component: ContributeAmountComponent;
  let fixture: ComponentFixture<ContributeAmountComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ContributeAmountComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContributeAmountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
